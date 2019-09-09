// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"container/list"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/muguangyi/ferry/network"
)

const (
	cOriginPort   int = 20000
	cMaxPortRange int = 49000
)

func newHub() *hub {
	return &hub{
		docks:       make(map[string]*list.List),
		assignPorts: make(map[string]int),
		blackPorts:  make(map[int]bool),
	}
}

type hub struct {
	socket           network.ISocket
	docksMutex       sync.Mutex
	docks            map[string]*list.List
	assignPortsMutex sync.Mutex
	assignPorts      map[string]int
	blackPorts       map[int]bool
}

type stub struct {
	addr  string
	peer  network.IPeer
	ready bool
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
				stubs, ok := h.docks[v]
				if !ok {
					stubs = list.New()
					h.docks[v] = stubs
				}
				stubs.PushBack(&stub{addr: addr, peer: peer, ready: false})
				h.docksMutex.Unlock()
			}

			resp := &packer{
				Id: cRegisterResponse,
				P: &protoRegisterResponse{
					Port: port,
				},
			}
			peer.Send(resp)
		}
	case cReady:
		{
			req := pack.P.(*protoReady)
			for _, slot := range req.Slots {
				h.docksMutex.Lock()
				stubs, ok := h.docks[slot]
				h.docksMutex.Unlock()
				if ok {
					for i := stubs.Front(); i != nil; i = i.Next() {
						stub := i.Value.(*stub)
						if stub.peer == peer {
							stub.ready = true
							break
						}
					}
				}
			}
		}
	case cQueryRequest:
		{
			go func() {
				req := pack.P.(*protoQueryRequest)
				for {
					h.docksMutex.Lock()
					docks, ok := h.docks[req.Slot]
					h.docksMutex.Unlock()
					if ok {
						// Loop from back to front, means the 'new' one will
						// be serve at first.
						for i := docks.Back(); i != nil; i = i.Prev() {
							stub := i.Value.(*stub)
							if stub.ready {
								h.respondQueryImme(peer, stub.addr)
								return
							}
						}
					}
				}
			}()
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

func (h *hub) respondQueryImme(peer network.IPeer, dockAddr string) {
	resp := &packer{
		Id: cQueryResponse,
		P: &protoQueryResponse{
			DockAddr: dockAddr,
		},
	}
	peer.Send(resp)
}

func pickIP(addr string) string {
	info := strings.Split(addr, ":")
	return info[0]
}
