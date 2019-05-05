// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	mock      bool                 = false
	listeners map[string]*listener = make(map[string]*listener)
	vport     int                  = 65536
)

func listen(network string, address string) (net.Listener, error) {
	if mock {
		address = formatAddr(address)
		listener := &listener{
			address: &addr{
				network: network,
				address: address,
			},
			chanconn: make(chan *conn, 1),
			conns:    make(map[string]*conn),
		}
		listeners[address] = listener

		return listener, nil
	} else {
		return net.Listen(network, address)
	}
}

func dial(network string, address string) (net.Conn, error) {
	if mock {
		var conn net.Conn = nil
		connected := make(chan bool, 1)
		go func() {
			address = formatAddr(address)
			timeout := 1 * time.Second
			for {
				listener := listeners[address]
				if nil != listener {
					vport += 1
					c := newConn()
					c.localAddr.network = network
					c.localAddr.address = fmt.Sprintf("0.0.0.0:%d", vport)
					listener.chanconn <- c

					conn = c
					connected <- true
					break
				} else if timeout <= 0 {
					connected <- false
					break
				}

				timeout -= time.Microsecond
				time.Sleep(time.Microsecond)
			}

		}()

		if succ := <-connected; succ {
			return conn, nil
		} else {
			return nil, errors.New("Can't connect server")
		}
	} else {
		return net.Dial(network, address)
	}
}

func reset() {
	mock = false
	listeners = make(map[string]*listener)
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

	c := newConn()
	c.localAddr = l.address
	c.remoteAddr = inconn.localAddr
	c.peer = inconn

	inconn.remoteAddr = l.address
	inconn.peer = c

	return c, nil
}

func (l *listener) Close() error {
	l.chanconn <- nil
	return nil
}

func (l *listener) Addr() net.Addr {
	return nil
}

func (l *listener) connect(conn *conn) {
	l.conns[conn.localAddr.String()] = conn
}

func newConn() *conn {
	c := new(conn)
	c.localAddr = new(addr)
	c.remoteAddr = new(addr)
	c.chanbuf = make(chan []byte, 1)

	return c
}

type conn struct {
	localAddr  *addr
	remoteAddr *addr
	peer       *conn
	chanbuf    chan []byte
}

func (c *conn) Read(b []byte) (n int, err error) {
	bytes := <-c.chanbuf
	return copy(b, bytes), nil
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
