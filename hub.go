// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"github.com/muguangyi/gounite/network"
)

func RunHub(port string) {
	addr := "0.0.0.0:" + port
	hub := network.NewSocket(addr, "json", &hubSink{})
	go hub.Listen()
}

type hubSink struct {
}

func (h *hubSink) OnConnected(p network.IPeer) {

}

func (h *hubSink) OnClosed(p network.IPeer) {

}

func (h *hubSink) OnPacket(p network.IPeer, obj interface{}) {

}
