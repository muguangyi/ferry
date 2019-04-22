// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/muguangyi/gounite/misc"
	"github.com/muguangyi/gounite/network"
)

type Hub struct {
	socket network.ISocket
	units  map[string][]string
	ports  map[string]int
}

func (h *Hub) Run(hubAddr string) {
	h.units = make(map[string][]string)
	h.ports = make(map[string]int)

	network.ExtendSerializer("gounite", newSerializer())

	h.socket = network.NewSocket(hubAddr, "gounite", h)
	go h.socket.Listen()
}

func (h *Hub) OnConnected(p network.IPeer) {

}

func (h *Hub) OnClosed(p network.IPeer) {

}

func (h *Hub) OnPacket(p network.IPeer, obj interface{}) {
	pack := obj.(jsonPack)
	switch pack.id {
	case REGISTER_REQUEST:
		{
			req := pack.p.(protoRegisterRequest)
			for _, v := range req.units {
				unions := h.units[v]
				unions = append(unions, p.LocalAddr().String())
			}

			addr := p.LocalAddr().String()
			port := h.ports[addr]
			if 0 == port {
				port = 50000
			}
			port += 1
			h.ports[addr] = port

			resp := &jsonPack{
				id: REGISTER_RESPONSE,
				p: &protoRegisterResponse{
					port: port,
				},
			}
			p.Send(resp)
		}
	case IMPORT_REQUEST:
		{
			req := pack.p.(protoImportRequest)
			set := misc.NewSet()
			for _, v := range req.units {
				unions := h.units[v]
				set.Add(unions[len(unions)-1])
			}

			slice := set.ToSlice()
			unions := make([]string, len(slice))
			for i, v := range slice {
				unions[i] = v.(string)
			}

			resp := &jsonPack{
				id: IMPORT_RESPONSE,
				p: &protoImportResponse{
					unions: unions,
				},
			}
			p.Send(resp)
		}
	case QUERY_REQUEST:
		{
			req := pack.p.(protoQueryRequest)
			unions := h.units[req.unit]
			resp := &jsonPack{
				id: QUERY_RESPONSE,
				p: &protoQueryResponse{
					unionAddr: unions[len(unions)-1],
				},
			}
			p.Send(resp)
		}
	}
}
