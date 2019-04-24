// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unite

import (
	"fmt"
	"strings"
	"sync"

	"github.com/muguangyi/gounite/misc"
	"github.com/muguangyi/gounite/network"
)

const (
	ORIGIN_PORT    int = 20000
	MAX_PORT_RANGE int = 49000
)

func newHub() *hub {
	return &hub{
		unitUnions:  make(map[string][]string),
		assignPorts: make(map[string]int),
		blackPorts:  make(map[int]bool),
	}
}

type hub struct {
	socket           network.ISocket
	unitUnionsMutex  sync.Mutex
	unitUnions       map[string][]string
	assignPortsMutex sync.Mutex
	assignPorts      map[string]int
	blackPorts       map[int]bool
}

func (h *hub) OnConnected(peer network.IPeer) {
}

func (h *hub) OnClosed(peer network.IPeer) {

}

func (h *hub) OnPacket(peer network.IPeer, obj interface{}) {
	pack := obj.(*packer)
	switch pack.Id {
	case REGISTER_REQUEST:
		{
			req := pack.P.(*protoRegisterRequest)

			addr := pickIP(peer.RemoteAddr().String())
			port := h.allocate(addr)

			addr = fmt.Sprintf("%s:%d", addr, port)
			for _, v := range req.Units {
				h.unitUnionsMutex.Lock()
				unions := h.unitUnions[v]
				if nil == unions {
					unions = make([]string, 0)
					h.unitUnions[v] = unions
				}
				h.unitUnions[v] = append(unions, addr)
				h.unitUnionsMutex.Unlock()
			}

			resp := &packer{
				Id: REGISTER_RESPONSE,
				P: &protoRegisterResponse{
					Port: port,
				},
			}
			peer.Send(resp)
		}
	case IMPORT_REQUEST:
		{
			go func() {
				req := pack.P.(*protoImportRequest)
				for {
					completed := true
					set := misc.NewSet()
					for _, v := range req.Units {
						h.unitUnionsMutex.Lock()
						unions := h.unitUnions[v]
						h.unitUnionsMutex.Unlock()
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

						resp := &packer{
							Id: IMPORT_RESPONSE,
							P: &protoImportResponse{
								Unions: unions,
							},
						}
						peer.Send(resp)

						return
					}
				}
			}()
		}
	case QUERY_REQUEST:
		{
			req := pack.P.(*protoQueryRequest)
			unions := h.unitUnions[req.Unit]
			resp := &packer{
				Id: QUERY_RESPONSE,
				P: &protoQueryResponse{
					UnionAddr: unions[len(unions)-1],
				},
			}
			peer.Send(resp)
		}
	}
}

func (h *hub) run(hubAddr string, blackPorts ...int) {
	for _, p := range blackPorts {
		h.blackPorts[p] = true
	}

	network.ExtendSerializer("gounite", newSerializer())

	h.socket = network.NewSocket(hubAddr, "gounite", h)
	go h.socket.Listen()
}

func (h *hub) allocate(addr string) int {
	h.assignPortsMutex.Lock()
	defer h.assignPortsMutex.Unlock()

	port := h.assignPorts[addr]
	for {
		if 0 == port {
			port = ORIGIN_PORT
		}

		port += 1

		if !h.blackPorts[port] {
			break
		}
	}

	if port > MAX_PORT_RANGE {
		panic("Out of port max range!")
	}

	h.assignPorts[addr] = port

	return port
}

func pickIP(addr string) string {
	info := strings.Split(addr, ":")
	return info[0]
}
