// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

import (
	"fmt"
	"strings"

	"github.com/muguangyi/gounite/misc"
	"github.com/muguangyi/gounite/network"
)

func NewHub() *Hub {
	return &Hub{
		unitUnions: make(map[string][]string),
		ports:      make(map[string]int),
	}
}

type Hub struct {
	socket     network.ISocket
	unitUnions map[string][]string
	ports      map[string]int
}

func (h *Hub) Run(hubAddr string) {
	network.ExtendSerializer("gounite", newSerializer())

	h.socket = network.NewSocket(hubAddr, "gounite", h)
	go h.socket.Listen()
}

func (h *Hub) OnConnected(p network.IPeer) {
}

func (h *Hub) OnClosed(p network.IPeer) {

}

func (h *Hub) OnPacket(p network.IPeer, obj interface{}) {
	pack := obj.(*jsonPack)
	switch pack.Id {
	case REGISTER_REQUEST:
		{
			req := pack.P.(*protoRegisterRequest)

			addr := strings.Split(p.RemoteAddr().String(), ":")[0]
			port := h.ports[addr]
			if 0 == port {
				port = 5000
			}
			port += 1
			h.ports[addr] = port

			addr = fmt.Sprintf("%s:%d", addr, port)
			for _, v := range req.Units {
				unions := h.unitUnions[v]
				if nil == unions {
					unions = make([]string, 0)
					h.unitUnions[v] = unions
				}
				h.unitUnions[v] = append(unions, addr)
			}

			resp := &jsonPack{
				Id: REGISTER_RESPONSE,
				P: &protoRegisterResponse{
					Port: port,
				},
			}
			p.Send(resp)
		}
	case IMPORT_REQUEST:
		{
			go func() {
				req := pack.P.(*protoImportRequest)
				for {
					completed := true
					set := misc.NewSet()
					for _, v := range req.Units {
						unions := h.unitUnions[v]
						if len(unions) > 0 {
							set.Add(unions[len(unions)-1])
						} else {
							completed = false
							break
						}
					}

					if completed {
						slice := set.ToSlice()
						unions := make([]string, len(slice))
						for i, v := range slice {
							unions[i] = v.(string)
						}

						resp := &jsonPack{
							Id: IMPORT_RESPONSE,
							P: &protoImportResponse{
								Unions: unions,
							},
						}
						p.Send(resp)

						return
					}
				}
			}()
		}
	case QUERY_REQUEST:
		{
			req := pack.P.(*protoQueryRequest)
			unions := h.unitUnions[req.Unit]
			resp := &jsonPack{
				Id: QUERY_RESPONSE,
				P: &protoQueryResponse{
					UnionAddr: unions[len(unions)-1],
				},
			}
			p.Send(resp)
		}
	}
}
