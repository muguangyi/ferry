// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"sync"
)

type ISerializer interface {
	Marshal(obj interface{}) []byte
	Unmarshal(data []byte) interface{}
	Slice(source []byte) int
}

type ISocket interface {
	Listen()
	Dial()
	Close()
}

type IPeer interface {
	IsSelf() bool
	Send(obj interface{})
}

type ISocketSink interface {
	OnConnected(p IPeer)
	OnClosed(p IPeer)
	OnPacket(p IPeer, obj interface{})
}

func NewSocket(addr string, serializer string, sink ISocketSink) ISocket {
	once.Do(func() {
		ExtendSerializer("txt", new(txtSerializer))
		ExtendSerializer("json", new(jsonSerializer))
	})
	s := new(socket)
	s.addr = addr
	s.serializer = serializers[serializer]
	s.sink = sink

	return s
}

func ExtendSerializer(name string, serializer ISerializer) {
	serializers[name] = serializer
}

var (
	serializers map[string]ISerializer = make(map[string]ISerializer)
	once        sync.Once
)
