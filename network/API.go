// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"net"
	"sync"
)

// ISerializer
type ISerializer interface {
	// Marshal, marshal object to []byte
	Marshal(obj interface{}) []byte

	// Unmarshal, unmarshal []byte to object
	Unmarshal(data []byte) interface{}

	// Slice, slice []byte to a packet with return packet length
	Slice(source []byte) int
}

// ISocket
type ISocket interface {
	// Listen, listen for all connecting peer
	Listen()

	// Dial, dial to target socket
	Dial()

	// Close, close the socket
	Close()

	// Send, send object to all connected peers
	Send(obj interface{})
}

// IPeer
type IPeer interface {
	// IsSelf, check if this peer is my own
	IsSelf() bool

	// LocalAddr, return LocalAddr type
	LocalAddr() net.Addr

	// RemoteAddr, return RemoteAddr type
	RemoteAddr() net.Addr

	// Send, send object to peer
	Send(obj interface{})
}

// ISocketSink, handle callback for socket events
type ISocketSink interface {
	// OnConnected, called when socket connected
	OnConnected(p IPeer)

	// OnClosed, called when socket disconnected
	OnClosed(p IPeer)

	// OnPacket, called when received a packet from socket
	OnPacket(p IPeer, obj interface{})
}

// NewSocket, new a socket with addr, serializer type and callback sink object
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

// ExtendSerializer, extend serializer type with name and handling object
func ExtendSerializer(name string, serializer ISerializer) {
	serializers[name] = serializer
}

var (
	serializers map[string]ISerializer = make(map[string]ISerializer)
	once        sync.Once
)
