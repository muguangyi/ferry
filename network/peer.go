// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

func newPeer(conn net.Conn, serializer ISerializer, sink ISocketSink, self bool) *peer {
	p := new(peer)
	p.conn = conn
	p.serializer = serializer
	p.sink = sink
	p.self = self
	p.sendPackets = make(chan interface{}, 100)
	p.readBlock = make([]byte, 1024*1024)
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
	readBlock   []byte
	recvBuffer  *bytes.Buffer
}

func (p *peer) IsSelf() bool {
	return p.self
}

func (p *peer) LocalAddr() net.Addr {
	return p.conn.LocalAddr()
}

func (p *peer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
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

		fmt.Println("====OnClosed:", p.LocalAddr().String(), "|", p.RemoteAddr().String())
		p.conn.Close()
	}()

	// Recv routine
	go func() {
		for {
			size, err := p.conn.Read(p.readBlock)
			if nil != err {
				break
			}

			fmt.Println(">>>>> Recv:", size)
			p.recvBuffer.Write(p.readBlock[:size])
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

		fmt.Println("====OnClosed:", p.LocalAddr().String(), "|", p.RemoteAddr().String())
		p.conn.Close()
	}()
}

func (p *peer) close() {
	p.Lock()
	defer p.Unlock()

	p.Send(nil)
}
