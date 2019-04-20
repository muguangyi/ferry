// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"github.com/muguangyi/gounite/internal"
	"github.com/muguangyi/gounite/network"
	"github.com/muguangyi/gounite/unit"
)

func Run(hubAddr string, units ...unit.IUnit) {
	var node = network.NewSocket(hubAddr, "json", &internal.UnionSink{})
	go node.Dial()

	// TODO:
	// 2. Get container port from hub
	// 3. Setup listen server on port
	// 4. Init all units
	// 5. Query dependent unit containers from hub if needed
	// 6. Connect to target containers
	// 7. Register to each other for units
}

func RunHub(hubAddr string) {
	hub := network.NewSocket(hubAddr, "json", &internal.HubSink{})
	go hub.Listen()
}
