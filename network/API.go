// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"net"
	"sync"
)

// ISerializer interface.
type ISerializer interface {
	// Marshal object to []byte.
	Marshal(obj interface{}) []byte

	// Unmarshal []byte to object.
	Unmarshal(data []byte) interface{}

	// Slice []byte to a packet and return the packet length.
	Slice(source []byte) int
}

// ISocket interface.
type ISocket interface {
	// Listen for all connecting peers.
	Listen()

	// Dial to target socket.
	Dial()

	// Close the socket.
	Close()

	// Send object to all connected peers.
	Send(obj interface{})
}

// IPeer interface.
type IPeer interface {
	// Check if this peer is my own.
	IsSelf() bool

	// Return LocalAddr type.
	LocalAddr() net.Addr

	// Return RemoteAddr type.
	RemoteAddr() net.Addr

	// Send object to peer.
	Send(obj interface{})
}

// ISocketSink interface to handle callback for socket events.
type ISocketSink interface {
	// Called when socket connected.
	OnConnected(p IPeer)

	// Called when socket disconnected.
	OnClosed(p IPeer)

	// Called when received a packet from socket.
	OnPacket(p IPeer, obj interface{})
}

// NewSocket create a socket with addr, serializer type and callback sink object.
func NewSocket(addr string, serializer string, sink ISocketSink) ISocket {
	once.Do(func() {
		ExtendSerializer("txt", new(txtSerializer))
	})
	s := new(socket)
	s.addr = addr
	s.serializer = serializers[serializer]
	s.sink = sink

	return s
}

// Mock enable or disable mockup network for testing
func Mock(enable bool) {
	mock = enable
}

// ExtendSerializer extend serializer type with name and handling object.
func ExtendSerializer(name string, serializer ISerializer) {
	serializers[name] = serializer
}

var (
	serializers map[string]ISerializer = make(map[string]ISerializer)
	once        sync.Once
)
