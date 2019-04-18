// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"fmt"
	"net"
)

type socket struct {
	addr  string
	peers []*peer
}

func (s *socket) Listen() {
	listener, err := net.Listen("tcp", s.addr)
	if nil != err {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if nil != err {
			fmt.Println(err)
		}

		p := newPeer(conn)
		p.run()
		s.peers = append(s.peers, p)
	}
}

func (s *socket) Dial() {
	conn, err := net.Dial("tcp", s.addr)
	if nil != err {
		fmt.Println(err)
		return
	}

	p := newPeer(conn)
	p.run()
	s.peers = append(s.peers, p)
}

func (s *socket) Close() {

}
