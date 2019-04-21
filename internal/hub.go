// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/muguangyi/gounite/network"
)

type Hub struct {
	socket   network.ISocket
	profiles map[string][]string
}

func (h *Hub) Run(hubAddr string) {
	h.socket = network.NewSocket(hubAddr, "json", &hubSink{})
	go h.socket.Listen()
}

type hubSink struct {
}

func (s *hubSink) OnConnected(p network.IPeer) {

}

func (s *hubSink) OnClosed(p network.IPeer) {

}

func (s *hubSink) OnPacket(p network.IPeer, obj interface{}) {

}
