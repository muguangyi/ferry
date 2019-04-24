// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

import (
	"fmt"

	"github.com/muguangyi/gounite/chancall"
	"github.com/muguangyi/gounite/network"
)

func NewUnion(name string, units ...IUnit) *Union {
	union := new(Union)
	union.name = name
	union.localUnits = make(map[string]*unit)
	union.remoteUnits = make(map[string]network.IPeer)
	union.dialUnions = make(map[string]bool)
	union.rpcs = make(map[int64]*rpc)

	for _, v := range units {
		u := v.(*unit)
		union.localUnits[u.id] = u
		u.union = union
	}

	return union
}

type Union struct {
	name        string
	localUnits  map[string]*unit
	remoteUnits map[string]network.IPeer
	dialUnions  map[string]bool
	rpcs        map[int64]*rpc
}

func (u *Union) Run(hubAddr string) {
	network.ExtendSerializer("gounite", newSerializer())

	var socket = network.NewSocket(hubAddr, "gounite", u)
	go socket.Dial()
}

func (u *Union) OnConnected(p network.IPeer) {
	go func() {
		req := &jsonPack{
			Id: REGISTER_REQUEST,
			P: &protoRegisterRequest{
				Units: u.collect(),
			},
		}
		p.Send(req)
	}()
}

func (u *Union) OnClosed(p network.IPeer) {
}

func (u *Union) OnPacket(p network.IPeer, obj interface{}) {
	pack := obj.(*jsonPack)
	switch pack.Id {
	case ERROR:
		{

		}
	case REGISTER_REQUEST:
		{
			req := pack.P.(*protoRegisterRequest)
			for _, v := range req.Units {
				u.remoteUnits[v] = p
			}

			addr := p.RemoteAddr().String()
			if !u.dialUnions[addr] {
				u.dialUnions[addr] = true

				u.tryStart()
			}
		}
	case REGISTER_RESPONSE:
		{
			resp := pack.P.(*protoRegisterResponse)
			listenAddr := fmt.Sprintf("127.0.0.1:%d", resp.Port)
			socket := network.NewSocket(listenAddr, "gounite", u)
			go socket.Listen()

			go func() {
				u.init()

				req := &jsonPack{
					Id: IMPORT_REQUEST,
					P: &protoImportRequest{
						Units: u.depends(),
					},
				}
				p.Send(req)
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
			node := network.NewSocket(resp.UnionAddr, "gounite", u)
			go node.Dial()
		}
	case RPC_REQUEST:
		{
			req := pack.P.(*protoRpcRequest)
			target := u.localUnits[req.UnitId]
			if nil != target {
				go func() {
					caller := chancall.NewCaller(target.callee)
					if req.WithResult {
						result, _ := caller.CallWithResult(req.Method, req.Args...)
						resp := &jsonPack{
							Id: RPC_RESPONSE,
							P: &protoRpcResponse{
								Index:  req.Index,
								UnitId: req.UnitId,
								Method: req.Method,
								Result: result,
							},
						}
						p.Send(resp)
					} else {
						caller.Call(req.Method, req.Args...)
					}
				}()
			}
		}
	case RPC_RESPONSE:
		{
			resp := pack.P.(*protoRpcResponse)
			rpc := u.rpcs[resp.Index]
			if nil != rpc {
				go func() {
					rpc.callback(resp.Result)
					delete(u.rpcs, resp.Index)
				}()
			}
		}
	}
}

func (u *Union) init() {
	for _, v := range u.localUnits {
		v.control.OnInit(v)
	}
}

func (u *Union) collect() []string {
	ids := make([]string, 0)
	for id := range u.localUnits {
		ids = append(ids, id)
	}

	return ids
}

func (u *Union) depends() []string {
	ids := make([]string, 0)
	for _, v := range u.localUnits {
		ids = append(ids, v.depends...)
	}

	return ids
}

func (u *Union) tryStart() {
	for _, v := range u.dialUnions {
		if !v {
			return
		}
	}

	go u.start()
}

func (u *Union) start() {
	for _, v := range u.localUnits {
		v.control.OnStart()
	}
}

func (u *Union) invoke(r *rpc) {
	u.rpcs[r.index] = r
}
