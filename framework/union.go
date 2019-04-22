// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

import (
	"github.com/muguangyi/gounite/network"
)

func NewUnion() *Union {
	return &Union{
		dialUnions: make(map[string]bool),
	}
}

type Union struct {
	dialUnions map[string]bool
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
			units: localCollect(),
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
				remoteUnits[v] = p
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

			localInit()

			req := &jsonPack{
				id: IMPORT_REQUEST,
				p: &protoImportRequest{
					units: localDepends(),
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
				localStart()
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
			// req := pack.p.(protoRpcRequest)
		}
	case RPC_RESPONSE:
		{

		}
	}
}

func (u *Union) tryStart() {
	for _, v := range u.dialUnions {
		if !v {
			return
		}
	}

	localStart()
}
