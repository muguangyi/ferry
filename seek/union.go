// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

import (
	"errors"
	"fmt"

	"github.com/muguangyi/seek/chancall"
	"github.com/muguangyi/seek/network"
)

func newUnion(name string, signalers ...ISignaler) *union {
	union := new(union)
	union.name = name
	union.localSignalers = make(map[string]*signaler)
	union.remoteSignalers = make(map[string]network.IPeer)
	union.dialUnions = make(map[string]bool)
	union.rpcs = make(map[int64]*rpc)

	for _, v := range signalers {
		s := v.(*signaler)
		union.localSignalers[s.callee.Name()] = s
		s.union = union
	}

	return union
}

type union struct {
	name            string
	localSignalers  map[string]*signaler
	remoteSignalers map[string]network.IPeer
	dialUnions      map[string]bool
	rpcs            map[int64]*rpc
}

func (u *union) OnConnected(peer network.IPeer) {
	go func() {
		req := &packer{
			Id: cREGISTER_REQUEST,
			P: &protoRegisterRequest{
				Signalers: u.collect(),
			},
		}
		peer.Send(req)
	}()
}

func (u *union) OnClosed(peer network.IPeer) {
}

func (u *union) OnPacket(peer network.IPeer, obj interface{}) {
	pack := obj.(*packer)
	switch pack.Id {
	case cERROR:
		{

		}
	case cREGISTER_REQUEST:
		{
			req := pack.P.(*protoRegisterRequest)
			for _, v := range req.Signalers {
				u.remoteSignalers[v] = peer
			}

			addr := peer.RemoteAddr().String()
			if !u.dialUnions[addr] {
				u.dialUnions[addr] = true

				u.tryStart()
			}
		}
	case cREGISTER_RESPONSE:
		{
			resp := pack.P.(*protoRegisterResponse)
			listenAddr := fmt.Sprintf("0.0.0.0:%d", resp.Port)
			socket := network.NewSocket(listenAddr, "seek", u)
			go socket.Listen()

			go func() {
				u.init()

				req := &packer{
					Id: cIMPORT_REQUEST,
					P: &protoImportRequest{
						Signalers: u.depends(),
					},
				}
				peer.Send(req)
			}()
		}
	case cIMPORT_RESPONSE:
		{
			resp := pack.P.(*protoImportResponse)
			if len(resp.Unions) > 0 {
				u.dialUnions = make(map[string]bool)
				for _, v := range resp.Unions {
					socket := network.NewSocket(v, "seek", u)
					go socket.Dial()

					u.dialUnions[v] = false
				}
			} else {
				go u.start()
			}
		}
	case cQUERY_RESPONSE:
		{
			resp := pack.P.(*protoQueryResponse)
			socket := network.NewSocket(resp.UnionAddr, "seek", u)
			go socket.Dial()
		}
	case cRPC_REQUEST:
		{
			req := pack.P.(*protoRpcRequest)
			target := u.localSignalers[req.SignalerId]
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
						Id: cRPC_RESPONSE,
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
	case cRPC_RESPONSE:
		{
			resp := pack.P.(*protoRpcResponse)
			rpc := u.rpcs[resp.Index]
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
					delete(u.rpcs, resp.Index)
				}()
			}
		}
	}
}

func (u *union) run(hubAddr string) {
	network.ExtendSerializer("seek", newSerializer())

	var socket = network.NewSocket(hubAddr, "seek", u)
	go socket.Dial()
}

func (u *union) init() {
	for _, v := range u.localSignalers {
		v.signal.OnInit(v)
	}
}

func (u *union) collect() []string {
	ids := make([]string, 0)
	for id, v := range u.localSignalers {
		if v.discoverable {
			ids = append(ids, id)
		}
	}

	return ids
}

func (u *union) depends() []string {
	ids := make([]string, 0)
	for _, v := range u.localSignalers {
		ids = append(ids, v.depends...)
	}

	return ids
}

func (u *union) tryStart() {
	for _, v := range u.dialUnions {
		if !v {
			return
		}
	}

	go u.start()
}

func (u *union) start() {
	for _, v := range u.localSignalers {
		v.signal.OnStart()
	}
}

func (u *union) invoke(rpc *rpc) {
	u.rpcs[rpc.index] = rpc
}

func (u *union) call(name string, method string, args ...interface{}) error {
	target := u.localSignalers[name]
	if nil != target {
		return chancall.NewCaller(target.callee).Call(method, args...)
	} else if p := u.remoteSignalers[name]; nil != p {
		rpc := newRpc()
		u.rpcs[rpc.index] = rpc
		return rpc.call(p, name, method, args...)
	}

	return fmt.Errorf("NO [%s] unit exist!", name)
}

func (u *union) callWithResult(name string, method string, args ...interface{}) ([]interface{}, error) {
	target := u.localSignalers[name]
	if nil != target {
		return chancall.NewCaller(target.callee).CallWithResult(method, args...)
	} else if p := u.remoteSignalers[name]; nil != p {
		rpc := newRpc()
		u.rpcs[rpc.index] = rpc
		return rpc.callWithResult(p, name, method, args...)
	}

	return nil, fmt.Errorf("NO [%s] unit exist!", name)
}
