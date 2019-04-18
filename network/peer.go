// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"net"
	"sync"
)

const (
	BUFFERSIZE int = 1024 * 1024
)

func newPeer(conn net.Conn) *peer {
	p := new(peer)
	p.conn = conn
	p.sendBuf = make(chan []byte)
	p.recvBuf = make([]byte, BUFFERSIZE)

	return p
}

type peer struct {
	sync.Mutex
	conn    net.Conn
	sendBuf chan []byte
	recvBuf []byte
}

func (p *peer) run() {
	go func() {
		for b := range p.sendBuf {
			if nil == b {
				break
			}

			_, err := p.conn.Write(b)
			if nil != err {
				break
			}
		}

		p.conn.Close()
	}()

	go func() {
		for {
			// count, err := p.conn.Read(p.recvBuf)
			// if nil != err {

			// }

		}
	}()
}

func (p *peer) write(b []byte) {
	p.Lock()
	defer p.Unlock()

	if len(p.sendBuf) == cap(p.sendBuf) {
		return
	}

	p.sendBuf <- b
}
