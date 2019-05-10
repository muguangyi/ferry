// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ship

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/muguangyi/ship/misc"
	"github.com/muguangyi/ship/network"
)

const (
	cOriginPort   int = 20000
	cMaxPortRange int = 49000
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

func (h *hub) Close() {
	h.socket.Close()
	h.socket = nil
	h.unitUnions = nil
	h.assignPorts = nil
	h.blackPorts = nil
}

func (h *hub) OnConnected(peer network.IPeer) {
}

func (h *hub) OnClosed(peer network.IPeer) {

}

func (h *hub) OnPacket(peer network.IPeer, obj interface{}) {
	pack := obj.(*packer)
	switch pack.Id {
	case cRegisterRequest:
		{
			req := pack.P.(*protoRegisterRequest)

			addr := pickIP(peer.RemoteAddr().String())
			port := h.allocate(addr)

			addr = fmt.Sprintf("%s:%d", addr, port)
			for _, v := range req.Signalers {
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
				Id: cRegisterResponse,
				P: &protoRegisterResponse{
					Port: port,
				},
			}
			peer.Send(resp)
		}
	case cImportRequest:
		{
			go func() {
				req := pack.P.(*protoImportRequest)
				for {
					completed := true
					set := misc.NewSet()
					for _, v := range req.Signalers {
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
							Id: cImportResponse,
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
	case cQueryRequest:
		{
			req := pack.P.(*protoQueryRequest)
			unions := h.unitUnions[req.Signaler]
			resp := &packer{
				Id: cQueryResponse,
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

	network.ExtendSerializer("seek", newSerializer())

	h.socket = network.NewSocket(hubAddr, "seek", h)
	h.socket.Listen()
}

func (h *hub) allocate(addr string) int {
	h.assignPortsMutex.Lock()
	defer h.assignPortsMutex.Unlock()

	port := h.assignPorts[addr]
	for {
		if 0 == port {
			port = cOriginPort
		}

		port += 1

		if !h.blackPorts[port] {
			break
		}
	}

	if port > cMaxPortRange {
		log.Fatal("Out of port max range!")
	}

	h.assignPorts[addr] = port

	return port
}

func pickIP(addr string) string {
	info := strings.Split(addr, ":")
	return info[0]
}
