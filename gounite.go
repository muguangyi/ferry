// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"net"

	"unit"
)

func Run(units ...unit.Unit) {
	for i := 0; i < len(units); i++ {

	}

	socket, _ := net.Listen("tcp", "0.0.0.0:17000")
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
