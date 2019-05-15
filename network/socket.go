// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"log"
	"net"
)

type socket struct {
	net        inet
	sink       ISocketSink
	addr       string
	listener   net.Listener
	serializer ISerializer
	peers      []*peer
}

func (s *socket) Listen() {
	var err error
	s.net, err = makeNet("tcp")
	if nil != err {
		log.Fatal(err)
		return
	}

	s.listener, err = s.net.Listen("tcp", s.addr)
	if nil != err {
		log.Fatal(err)
		return
	}

	go func() {
		for {
			conn, err := s.listener.Accept()
			if nil != err {
				log.Println(err)
				return
			}

			peer := newPeer(conn, s.serializer, s.sink, false)
			s.peers = append(s.peers, peer)

			if nil != s.sink {
				s.sink.OnConnected(peer)
			}

			peer.run()
		}
	}()
}

func (s *socket) Dial() {
	if nil == s.sink {
		log.Fatal("Please call Init first!")
		return
	}

	var err error
	s.net, err = makeNet("tcp")
	if nil != err {
		log.Fatal(err)
		return
	}

	conn, err := s.net.Dial("tcp", s.addr)
	if nil != err {
		log.Fatal(err)
		return
	}

	peer := newPeer(conn, s.serializer, s.sink, true)
	s.peers = append(s.peers, peer)

	if nil != s.sink {
		s.sink.OnConnected(peer)
	}

	peer.run()
}

func (s *socket) Close() {
	// NOTE: Only close listener, but not set to nil
	if nil != s.listener {
		s.listener.Close()
	}

	for _, peer := range s.peers {
		peer.close()
	}
	s.peers = nil
}

func (s *socket) Send(obj interface{}) {
	for _, peer := range s.peers {
		peer.Send(obj)
	}
}
