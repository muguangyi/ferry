// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"bytes"
	"net"
	"sync"
)

const (
	BUFFERSIZE int = 1024 * 1024
)

var (
	readBlock []byte = make([]byte, 1024*10)
)

func newPeer(conn net.Conn, serializer ISerializer) *peer {
	p := new(peer)
	p.conn = conn
	p.serializer = serializer
	p.sendPacket = make(chan *packet, 10)
	p.recvBuffer = new(bytes.Buffer)

	return p
}

type peer struct {
	sync.Mutex
	conn       net.Conn
	serializer ISerializer
	sendPacket chan *packet
	recvBuffer *bytes.Buffer
}

type packet struct {
	target interface{}
}

func (p *peer) run() {
	// Send routine
	go func() {
		for {
			err := p.send(<-p.sendPacket)
			if nil != err {
				break
			}
		}

		p.conn.Close()
	}()

	// Recv routine
	go func() {
		for {
			size, err := p.conn.Read(readBlock)
			if nil != err {

			}

			p.recvBuffer.Write(readBlock[:size])
			length := p.serializer.Slice(p.recvBuffer.Bytes())
			if length > 0 {
				block := make([]byte, length)
				n, err := p.recvBuffer.Read(block)
				if nil != err || n != length {
					p.serializer.Unmarshal(block)
				}
			}
		}
	}()
}

func (p *peer) write(obj interface{}) {
	p.Lock()
	defer p.Unlock()

	if len(p.sendPacket) == cap(p.sendPacket) {
		return
	}

	p.sendPacket <- &packet{target: obj}
}

func (p *peer) send(packet *packet) (err error) {
	data := p.serializer.Marshal(packet)
	_, err = p.conn.Write(data)
	return
}
