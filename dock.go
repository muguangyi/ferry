// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"errors"
	"fmt"
	"sync"

	"github.com/muguangyi/ferry/chancall"
	"github.com/muguangyi/ferry/network"
)

func newDock(name string, sandboxes ...ISandbox) *dock {
	d := new(dock)
	d.name = name
	d.sockets = make([]network.ISocket, 0)
	d.sandboxes = make(map[string]*sandbox)
	d.remoteSandboxes = make(map[string]network.IPeer)
	d.dialDocks = make(map[string]bool)
	d.rpcs = make(map[int64]*rpc)

	for _, v := range sandboxes {
		s := v.(*sandbox)
		d.sandboxes[s.callee.Name()] = s
		s.dock = d
	}

	return d
}

type dock struct {
	name                 string
	sockets              []network.ISocket
	sandboxes            map[string]*sandbox
	remoteSandboxesMutex sync.Mutex
	remoteSandboxes      map[string]network.IPeer
	dialDocks            map[string]bool
	rpcs                 map[int64]*rpc
}

func (d *dock) Close() {
	for _, v := range d.sandboxes {
		v.feature.OnDestroy()
	}
	d.sandboxes = nil

	for i := len(d.sockets) - 1; i >= 0; i-- {
		d.sockets[i].Close()
	}
	d.sockets = d.sockets[:0]
}

func (u *dock) OnConnected(peer network.IPeer) {
	go func() {
		req := &packer{
			Id: cRegisterRequest,
			P: &protoRegisterRequest{
				Signalers: u.collect(),
			},
		}
		peer.Send(req)
	}()
}

func (d *dock) OnClosed(peer network.IPeer) {
}

func (d *dock) OnPacket(peer network.IPeer, obj interface{}) {
	pack := obj.(*packer)
	switch pack.Id {
	case cError:
		{

		}
	case cRegisterRequest:
		{
			req := pack.P.(*protoRegisterRequest)
			for _, v := range req.Signalers {
				d.remoteSandboxesMutex.Lock()
				d.remoteSandboxes[v] = peer
				d.remoteSandboxesMutex.Unlock()
			}

			addr := peer.RemoteAddr().String()
			if !d.dialDocks[addr] {
				d.dialDocks[addr] = true

				d.tryStart()
			}
		}
	case cRegisterResponse:
		{
			resp := pack.P.(*protoRegisterResponse)
			listenAddr := fmt.Sprintf("0.0.0.0:%d", resp.Port)
			socket := network.NewSocket(listenAddr, "seek", d)
			socket.Listen()
			d.sockets = append(d.sockets, socket)

			go func() {
				d.init()

				req := &packer{
					Id: cImportRequest,
					P: &protoImportRequest{
						Signalers: d.depends(),
					},
				}
				peer.Send(req)
			}()
		}
	case cImportResponse:
		{
			resp := pack.P.(*protoImportResponse)
			if len(resp.Unions) > 0 {
				d.dialDocks = make(map[string]bool)
				for _, v := range resp.Unions {
					socket := network.NewSocket(v, "seek", d)
					socket.Dial()
					d.sockets = append(d.sockets, socket)

					d.dialDocks[v] = false
				}
			} else {
				go d.start()
			}
		}
	case cQueryResponse:
		{
			resp := pack.P.(*protoQueryResponse)
			socket := network.NewSocket(resp.UnionAddr, "seek", d)
			socket.Dial()
			d.sockets = append(d.sockets, socket)
		}
	case cRpcRequest:
		{
			req := pack.P.(*protoRpcRequest)
			target := d.sandboxes[req.SignalerId]
			if nil != target {
				go func() {
					caller := chancall.NewCaller(target.callee)
					var result []interface{}
					var err error
					if req.WithResult {
						result, err = caller.CallWithResult(req.Method, req.Args...)
					} else {
						err = caller.Call(req.Method, req.Args...)
					}

					resp := &packer{
						Id: cRpcResponse,
						P: &protoRpcResponse{
							Index:      req.Index,
							SignalerId: req.SignalerId,
							Method:     req.Method,
							Result:     result,
							Err: func() string {
								if nil != err {
									return err.Error()
								}

								return ""
							}(),
						},
					}
					peer.Send(resp)
				}()
			}
		}
	case cRpcResponse:
		{
			resp := pack.P.(*protoRpcResponse)
			rpc := d.rpcs[resp.Index]
			if nil != rpc {
				go func() {
					rpc.callback(&ret{
						result: resp.Result,
						err: func() error {
							if "" != resp.Err {
								return errors.New(resp.Err)
							}

							return nil
						}(),
					})
					delete(d.rpcs, resp.Index)
				}()
			}
		}
	}
}

func (d *dock) run(hubAddr string) {
	network.ExtendSerializer("seek", newSerializer())

	var socket = network.NewSocket(hubAddr, "seek", d)
	socket.Dial()
	d.sockets = append(d.sockets, socket)
}

func (d *dock) init() {
	for _, v := range d.sandboxes {
		v.feature.OnInit(v)
	}
}

func (d *dock) collect() []string {
	ids := make([]string, 0)
	for id, v := range d.sandboxes {
		if v.discoverable {
			ids = append(ids, id)
		}
	}

	return ids
}

func (d *dock) depends() []string {
	ids := make([]string, 0)
	for _, v := range d.sandboxes {
		ids = append(ids, v.depends...)
	}

	return ids
}

func (d *dock) tryStart() {
	for _, v := range d.dialDocks {
		if !v {
			return
		}
	}

	go d.start()
}

func (d *dock) start() {
	for _, v := range d.sandboxes {
		v.feature.OnStart()
	}
}

func (d *dock) call(name string, method string, args ...interface{}) error {
	target := d.sandboxes[name]
	if nil != target {
		return chancall.NewCaller(target.callee).Call(method, args...)
	} else if p := d.queryRemoteSignaler(name); nil != p {
		rpc := newRpc()
		d.rpcs[rpc.index] = rpc
		return rpc.call(p, name, method, args...)
	}

	return fmt.Errorf("NO [%s] unit exist!", name)
}

func (d *dock) callWithResult(name string, method string, args ...interface{}) ([]interface{}, error) {
	target := d.sandboxes[name]
	if nil != target {
		return chancall.NewCaller(target.callee).CallWithResult(method, args...)
	} else if p := d.queryRemoteSignaler(name); nil != p {
		rpc := newRpc()
		d.rpcs[rpc.index] = rpc
		return rpc.callWithResult(p, name, method, args...)
	}

	return nil, fmt.Errorf("NO [%s] unit exist!", name)
}

func (d *dock) queryRemoteSignaler(name string) network.IPeer {
	d.remoteSandboxesMutex.Lock()
	defer d.remoteSandboxesMutex.Unlock()

	return d.remoteSandboxes[name]
}
