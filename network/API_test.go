// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/muguangyi/ferry/network"
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
	c.wg.Done()
}

func (c *clientSink) OnClosed(p network.IPeer) {

}

func (c *clientSink) OnPacket(p network.IPeer, obj interface{}) {
	log.Println(obj)
	c.wg.Done()
}

func TestFlow(t *testing.T) {
	network.Mock("tcp")

	var wg sync.WaitGroup

	wg.Add(4)
	server := network.NewSocket("127.0.0.1:55555", "txt", &serverSink{wg: &wg})
	server.Listen()

	client := network.NewSocket("127.0.0.1:55555", "txt", &clientSink{wg: &wg})
	client.Dial()

	wg.Wait()
	client.Close()
	server.Close()
}
