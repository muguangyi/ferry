// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"fmt"
	"net"

	"github.com/muguangyi/gounite/unit"
)

func Run(units ...unit.Unit) {
	for i := 0; i < len(units); i++ {

	}

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
