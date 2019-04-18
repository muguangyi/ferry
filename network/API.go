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

func NewSocket(addr string, serializer string) ISocket {
	once.Do(func() {
		serializers["frame"] = new(frameSerializer)
	})
	s := new(socket)
	s.addr = addr
	s.serializer = serializers[serializer]

	return s
}

var (
	serializers map[string]ISerializer = make(map[string]ISerializer)
	once        sync.Once
)
