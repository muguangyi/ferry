// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"github.com/muguangyi/gounite/network"
	"github.com/muguangyi/gounite/unit"
)

func Run(hubAddr string, units ...unit.IUnit) {
	var node = network.NewSocket(hubAddr, "json", &gouniteSink{})
	go node.Dial()

	// TODO:
	// 2. Get container port from hub
	// 3. Setup listen server on port
	// 4. Init all units
	// 5. Query dependent unit containers from hub if needed
	// 6. Connect to target containers
	// 7. Register to each other for units
}

type gouniteSink struct {
}

func (g *gouniteSink) OnConnected(p network.IPeer) {
	if p.IsSelf() {
		// TODO:
		// 1.Register units to hub
	}
}

func (g *gouniteSink) OnClosed(p network.IPeer) {

}

func (g *gouniteSink) OnPacket(p network.IPeer, obj interface{}) {

}
