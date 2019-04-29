// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/muguangyi/seek/network"
)

type serverSink struct {
	wg *sync.WaitGroup
}

func (s *serverSink) OnConnected(p network.IPeer) {
	fmt.Println("One client connected!")
	p.Send("Hello\000 World!\n\000")
	s.wg.Done()
}

func (s *serverSink) OnClosed(p network.IPeer) {

}

func (s *serverSink) OnPacket(p network.IPeer, obj interface{}) {

}

type clientSink struct {
	wg *sync.WaitGroup
}

func (c *clientSink) OnConnected(p network.IPeer) {
	fmt.Println("Client connected!")
}

func (c *clientSink) OnClosed(p network.IPeer) {

}

func (c *clientSink) OnPacket(p network.IPeer, obj interface{}) {
	fmt.Print(obj)
	c.wg.Done()
}

func Test(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	server := network.NewSocket("0.0.0.0:55555", "txt", &serverSink{wg: &wg})
	go server.Listen()

	wg.Add(2)
	client := network.NewSocket("127.0.0.1:55555", "txt", &clientSink{wg: &wg})
	go client.Dial()

	wg.Wait()
	client.Close()
	server.Close()
}
