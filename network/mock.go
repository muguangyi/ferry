// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func makeNetMock(network string) inet {
	return new(netMockTcp)
}

var (
	listenersMutex sync.Mutex
	listeners      map[string]*listener = make(map[string]*listener)
	vport          int                  = 65536
)

type netMockTcp struct {
}

func (n *netMockTcp) Listen(network string, address string) (net.Listener, error) {
	fmt.Println("Listener started.")

	address = formatAddr(address)
	listener := &listener{
		address: &addr{
			network: network,
			address: address,
		},
		chanconn: make(chan *conn, 1),
		conns:    make(map[string]*conn),
	}
	listenersMutex.Lock()
	listeners[address] = listener
	listenersMutex.Unlock()

	return listener, nil
}

func (n *netMockTcp) Dial(network string, address string) (net.Conn, error) {
	c := newConn()
	connected := make(chan bool, 1)
	go func() {
		address = formatAddr(address)
		timeout := 1 * time.Second
		for {
			if 0 == c.status {
				listener := listeners[address]
				if nil != listener {
					vport += 1
					c.status = 1
					c.localAddr.network = network
					c.localAddr.address = fmt.Sprintf("0.0.0.0:%d", vport)
					c.remoteAddr = listener.address
					listener.chanconn <- c
				}
			} else if 2 == c.status {
				break
			} else if timeout <= 0 {
				c.status = -1
				break
			}

			timeout -= time.Microsecond
			time.Sleep(time.Microsecond)
		}

		connected <- (2 == c.status)
	}()

	succ := <-connected
	if succ {
		return c, nil
	} else {
		return nil, errors.New("Can't connect server")
	}
}

func reset() {
	listenersMutex.Lock()
	listeners = make(map[string]*listener)
	listenersMutex.Unlock()
	vport = 65535
}

func formatAddr(addr string) string {
	info := strings.Split(addr, ":")
	return "0.0.0.0:" + info[len(info)-1]
}

type listener struct {
	address  *addr
	chanconn chan *conn
	conns    map[string]*conn
}

func (l *listener) Accept() (net.Conn, error) {
	inconn := <-l.chanconn
	if nil == inconn {
		return nil, fmt.Errorf("Listener closed!")
	}

	fmt.Println("One connection is comming...")

	c := newConn()
	c.localAddr = l.address
	c.remoteAddr = inconn.localAddr
	c.peer = inconn

	inconn.peer = c
	inconn.status = 2

	l.conns[inconn.localAddr.String()] = inconn

	return c, nil
}

func (l *listener) Close() error {
	l.chanconn <- nil
	l.conns = nil
	fmt.Println("Listener closing...")
	return nil
}

func (l *listener) Addr() net.Addr {
	return l.address
}

func newConn() *conn {
	c := new(conn)
	c.status = 0
	c.localAddr = new(addr)
	c.remoteAddr = new(addr)
	c.chanbuf = make(chan []byte, 1)

	return c
}

type conn struct {
	status     int
	localAddr  *addr
	remoteAddr *addr
	peer       *conn
	chanbuf    chan []byte
}

func (c *conn) Read(b []byte) (n int, err error) {
	bytes := <-c.chanbuf
	if nil == bytes {
		return 0, fmt.Errorf("Target peer is closed!")
	} else {
		return copy(b, bytes), nil
	}
}

func (c *conn) Write(b []byte) (n int, err error) {
	if nil != c.peer {
		c.peer.chanbuf <- b
		return len(b), nil
	} else {
		return 0, fmt.Errorf("Conn is nil!")
	}
}

func (c *conn) Close() error {
	c.Write(nil)
	c.peer = nil
	return nil
}

func (c *conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return nil
}

type addr struct {
	network string
	address string
}

func (a *addr) Network() string {
	return a.network
}

func (a *addr) String() string {
	return a.address
}
