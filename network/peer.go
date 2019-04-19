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

func newPeer(conn net.Conn, serializer ISerializer, sink ISocketSink, self bool) *peer {
	p := new(peer)
	p.conn = conn
	p.serializer = serializer
	p.sink = sink
	p.self = self
	p.sendPackets = make(chan interface{}, 10)
	p.recvBuffer = new(bytes.Buffer)

	return p
}

type peer struct {
	sync.Mutex
	conn        net.Conn
	serializer  ISerializer
	sink        ISocketSink
	self        bool
	sendPackets chan interface{}
	recvBuffer  *bytes.Buffer
}

func (p *peer) IsSelf() bool {
	return p.self
}

func (p *peer) Send(obj interface{}) {
	if len(p.sendPackets) == cap(p.sendPackets) {
		return
	}

	p.sendPackets <- obj
}

func (p *peer) run() {
	// Send routine
	go func() {
		for {
			o := <-p.sendPackets
			if nil == o {
				break
			}

			data := p.serializer.Marshal(o)
			_, err := p.conn.Write(data)
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
			for length := p.serializer.Slice(p.recvBuffer.Bytes()); length > 0; length = p.serializer.Slice(p.recvBuffer.Bytes()) {
				slice := make([]byte, length)
				n, err := p.recvBuffer.Read(slice)
				if nil != err || n != length {
				}

				obj := p.serializer.Unmarshal(slice)
				if nil != p.sink {
					p.sink.OnPacket(p, obj)
				}
			}
		}
	}()
}

func (p *peer) close() {
	p.Lock()
	defer p.Unlock()

	p.Send(nil)
}
