// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

import (
	"github.com/muguangyi/gounite/chancall"
	"github.com/muguangyi/gounite/network"
)

func NewUnion(units ...IUnit) *Union {
	union := new(Union)
	union.remoteUnits = make(map[string]network.IPeer)
	union.dialUnions = make(map[string]bool)

	for _, v := range units {
		u := v.(*unit)
		union.localUnits[u.id] = u
		u.union = union
	}

	return union
}

type Union struct {
	localUnits  map[string]*unit
	remoteUnits map[string]network.IPeer
	dialUnions  map[string]bool
	rpcs        map[int64]*rpc
}

func (u *Union) Run(hubAddr string) {
	network.ExtendSerializer("gounite", newSerializer())

	var node = network.NewSocket(hubAddr, "gounite", u)
	go node.Dial()
}

func (u *Union) OnConnected(p network.IPeer) {
	req := &jsonPack{
		id: REGISTER_REQUEST,
		p: &protoRegisterRequest{
			units: u.collect(),
		},
	}
	p.Send(req)
}

func (u *Union) OnClosed(p network.IPeer) {

}

func (u *Union) OnPacket(p network.IPeer, obj interface{}) {
	pack := obj.(jsonPack)
	switch pack.id {
	case REGISTER_REQUEST:
		{
			req := pack.p.(protoRegisterRequest)
			for _, v := range req.units {
				u.remoteUnits[v] = p
			}

			addr := p.LocalAddr().String()
			if !u.dialUnions[addr] {
				u.dialUnions[addr] = true

				u.tryStart()
			}
		}
	case REGISTER_RESPONSE:
		{
			resp := pack.p.(protoRegisterResponse)
			listenAddr := "0.0.0.0:" + string(resp.port)
			socket := network.NewSocket(listenAddr, "gounite", u)
			go socket.Listen()

			u.init()

			req := &jsonPack{
				id: IMPORT_REQUEST,
				p: &protoImportRequest{
					units: u.depends(),
				},
			}
			p.Send(req)
		}
	case IMPORT_RESPONSE:
		{
			resp := pack.p.(protoImportResponse)
			if len(resp.unions) > 0 {
				u.dialUnions = make(map[string]bool)
				for _, v := range resp.unions {
					socket := network.NewSocket(v, "gounite", u)
					go socket.Dial()

					u.dialUnions[v] = false
				}
			} else {
				u.start()
			}
		}
	case QUERY_RESPONSE:
		{
			resp := pack.p.(protoQueryResponse)
			node := network.NewSocket(resp.unionAddr, "gounite", u)
			go node.Dial()
		}
	case RPC_REQUEST:
		{
			req := pack.p.(protoRpcRequest)
			target := u.localUnits[req.unitName]
			if nil != target {
				caller := chancall.NewCaller(target.callee)
				if req.withResult {
					result, err := caller.CallWithResult(req.method, req.args)
					resp := &jsonPack{
						id: RPC_RESPONSE,
						p: &protoRpcResponse{
							index:    req.index,
							unitName: req.unitName,
							method:   req.method,
							result:   result,
							err:      err,
						},
					}
					p.Send(resp)
				} else {
					caller.Call(req.method, req.args)
				}
			}
		}
	case RPC_RESPONSE:
		{
			resp := pack.p.(protoRpcResponse)
			rpc := u.rpcs[resp.index]
			if nil != rpc {
				rpc.callback(resp.result, resp.err)
				delete(u.rpcs, resp.index)
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

	u.start()
}

func (u *Union) start() {
	for _, v := range u.localUnits {
		v.control.OnStart()
	}
}

func (u *Union) invoke(r *rpc) {
	u.rpcs[r.index] = r
}
