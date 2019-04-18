// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"fmt"
	"net"

	"github.com/muguangyi/gounite/unit"
)

func Run(units ...unit.IUnit) {
	// TODO:
	// 1. Connect to hub and register current units to hub
	// 2. Get container port from hub
	// 3. Setup listen server on port
	// 4. Init all units
	// 5. Query dependent unit containers from hub if needed
	// 6. Connect to target containers
	// 7. Register to each other for units

	socket, error := net.Listen("tcp", "0.0.0.0:17000")
	if nil != error {
		fmt.Println(error)
	}

	defer socket.Close()

	for {
		conn, error := socket.Accept()
		if nil != error {
			continue
		}

		go recv(conn)
	}
}

func recv(conn net.Conn) {

}
