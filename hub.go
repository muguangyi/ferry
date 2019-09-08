// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/muguangyi/ferry/misc"
	"github.com/muguangyi/ferry/network"
)

const (
	cOriginPort   int = 20000
	cMaxPortRange int = 49000
)

func newHub() *hub {
	return &hub{
		docks:       make(map[string][]string),
		assignPorts: make(map[string]int),
		blackPorts:  make(map[int]bool),
	}
}

type hub struct {
	socket           network.ISocket
	docksMutex       sync.Mutex
	docks            map[string][]string
	assignPortsMutex sync.Mutex
	assignPorts      map[string]int
	blackPorts       map[int]bool
}

func (h *hub) Close() {
	h.socket.Close()
	h.socket = nil
	h.docks = nil
	h.assignPorts = nil
	h.blackPorts = nil
}

func (h *hub) OnConnected(peer network.IPeer) {
	log.Printf("[%s] dock is comming to hub [%s]...", peer.RemoteAddr(), peer.LocalAddr())
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
			for _, v := range req.Slots {
				h.docksMutex.Lock()
				docks := h.docks[v]
				if nil == docks {
					docks = make([]string, 0)
					h.docks[v] = docks
				}
				h.docks[v] = append(docks, addr)
				h.docksMutex.Unlock()
			}

			resp := &packer{
				Id: cHubRegisterResponse,
				P: &protoHubRegisterResponse{
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
					for _, v := range req.Slots {
						h.docksMutex.Lock()
						docks := h.docks[v]
						h.docksMutex.Unlock()
						if len(docks) > 0 {
							set.Add(docks[len(docks)-1])
						} else {
							completed = false
							break
						}
					}

					if completed {
						slice := set.ToSlice()
						docks := make([]string, len(slice))
						for i, v := range slice {
							docks[i] = v.(string)
						}

						resp := &packer{
							Id: cImportResponse,
							P: &protoImportResponse{
								Docks: docks,
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
			docks := h.docks[req.Slot]
			resp := &packer{
				Id: cQueryResponse,
				P: &protoQueryResponse{
					DockAddr: docks[len(docks)-1],
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
