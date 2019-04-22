// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/muguangyi/gounite/network"
	"github.com/muguangyi/gounite/unit"
)

func NewUnion() *Union {
	return &Union{
		units: make(map[string]network.IPeer),
	}
}

type Union struct {
	units map[string]network.IPeer
}

func (u *Union) Run(hubAddr string) {
	network.ExtendSerializer("gounite", newSerializer())

	var node = network.NewSocket(hubAddr, "gounite", u)
	go node.Dial()

	// TODO:
	// 5. Query dependent unit containers from hub if needed
	// 6. Connect to target containers
	// 7. Register to each other for units
}

func (u *Union) OnConnected(p network.IPeer) {
	req := &jsonPack{
		id: REGISTER_REQUEST,
		p: &protoRegisterRequest{
			units: unit.Collect(),
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
				u.units[v] = p
			}
		}
	case REGISTER_RESPONSE:
		{
			resp := pack.p.(protoRegisterResponse)
			listenAddr := "0.0.0.0:" + string(resp.port)
			socket := network.NewSocket(listenAddr, "gounite", u)
			go socket.Listen()

			unit.Init()

			req := &jsonPack{
				id: IMPORT_REQUEST,
				p: &protoImportRequest{
					units: unit.Depends(),
				},
			}
			p.Send(req)
		}
	case IMPORT_RESPONSE:
		{
			resp := pack.p.(protoImportResponse)
			if len(resp.unions) > 0 {
				for _, v := range resp.unions {
					socket := network.NewSocket(v, "gounite", u)
					go socket.Dial()
				}
			} else {
				// TODO:
				// Start union
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
