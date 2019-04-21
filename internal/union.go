// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/muguangyi/gounite/network"
	"github.com/muguangyi/gounite/unit"
)

type Union struct {
}

func (u *Union) Run(hubAddr string, units ...unit.IUnit) {
	var node = network.NewSocket(hubAddr, "json", &unionSink{})
	go node.Dial()

	// TODO:
	// 2. Get container port from hub
	// 3. Setup listen server on port
	// 4. Init all units
	// 5. Query dependent unit containers from hub if needed
	// 6. Connect to target containers
	// 7. Register to each other for units
}

type unionSink struct {
}

func (s *unionSink) OnConnected(p network.IPeer) {
	if p.IsSelf() {
		// TODO:
		// 1.Register units to hub
	}
}

func (s *unionSink) OnClosed(p network.IPeer) {

}

func (s *unionSink) OnPacket(p network.IPeer, obj interface{}) {

}
