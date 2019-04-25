// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unite

import (
	"errors"
	"fmt"

	"github.com/muguangyi/unite/chancall"
	"github.com/muguangyi/unite/network"
)

func newUnion(name string, units ...IUnit) *union {
	union := new(union)
	union.name = name
	union.localUnits = make(map[string]*unit)
	union.remoteUnits = make(map[string]network.IPeer)
	union.dialUnions = make(map[string]bool)
	union.rpcs = make(map[int64]*rpc)

	for _, v := range units {
		u := v.(*unit)
		union.localUnits[u.callee.Name()] = u
		u.union = union
	}

	return union
}

type union struct {
	name        string
	localUnits  map[string]*unit
	remoteUnits map[string]network.IPeer
	dialUnions  map[string]bool
	rpcs        map[int64]*rpc
}

func (u *union) OnConnected(peer network.IPeer) {
	go func() {
		req := &packer{
			Id: REGISTER_REQUEST,
			P: &protoRegisterRequest{
				Units: u.collect(),
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
	case ERROR:
		{

		}
	case REGISTER_REQUEST:
		{
			req := pack.P.(*protoRegisterRequest)
			for _, v := range req.Units {
				u.remoteUnits[v] = peer
			}

			addr := peer.RemoteAddr().String()
			if !u.dialUnions[addr] {
				u.dialUnions[addr] = true

				u.tryStart()
			}
		}
	case REGISTER_RESPONSE:
		{
			resp := pack.P.(*protoRegisterResponse)
			listenAddr := fmt.Sprintf("0.0.0.0:%d", resp.Port)
			socket := network.NewSocket(listenAddr, "gounite", u)
			go socket.Listen()

			go func() {
				u.init()

				req := &packer{
					Id: IMPORT_REQUEST,
					P: &protoImportRequest{
						Units: u.depends(),
					},
				}
				peer.Send(req)
			}()
		}
	case IMPORT_RESPONSE:
		{
			resp := pack.P.(*protoImportResponse)
			if len(resp.Unions) > 0 {
				u.dialUnions = make(map[string]bool)
				for _, v := range resp.Unions {
					socket := network.NewSocket(v, "gounite", u)
					go socket.Dial()

					u.dialUnions[v] = false
				}
			} else {
				go u.start()
			}
		}
	case QUERY_RESPONSE:
		{
			resp := pack.P.(*protoQueryResponse)
			socket := network.NewSocket(resp.UnionAddr, "gounite", u)
			go socket.Dial()
		}
	case RPC_REQUEST:
		{
			req := pack.P.(*protoRpcRequest)
			target := u.localUnits[req.UnitId]
			if nil != target {
				go func() {
					caller := chancall.NewCaller(target.callee)
					var result interface{}
					var err error
					if req.WithResult {
						result, err = caller.CallWithResult(req.Method, req.Args...)
					} else {
						err = caller.Call(req.Method, req.Args...)
					}

					resp := &packer{
						Id: RPC_RESPONSE,
						P: &protoRpcResponse{
							Index:  req.Index,
							UnitId: req.UnitId,
							Method: req.Method,
							Result: result,
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
	case RPC_RESPONSE:
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
	network.ExtendSerializer("gounite", newSerializer())

	var socket = network.NewSocket(hubAddr, "gounite", u)
	go socket.Dial()
}

func (u *union) init() {
	for _, v := range u.localUnits {
		v.control.OnInit(v)
	}
}

func (u *union) collect() []string {
	ids := make([]string, 0)
	for id := range u.localUnits {
		ids = append(ids, id)
	}

	return ids
}

func (u *union) depends() []string {
	ids := make([]string, 0)
	for _, v := range u.localUnits {
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
	for _, v := range u.localUnits {
		v.control.OnStart()
	}
}

func (u *union) invoke(rpc *rpc) {
	u.rpcs[rpc.index] = rpc
}
