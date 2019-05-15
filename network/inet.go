// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"fmt"
	"net"
)

type inet interface {
	Listen(network string, address string) (net.Listener, error)
	Dial(network string, address string) (net.Conn, error)
}

func bind(network string, maker netmaker) error {
	if "" == network || nil == maker {
		return fmt.Errorf("Input parameter is incorrect!")
	}

	makers[network] = maker
	return nil
}

type netmaker func(string) inet

func makeNet(network string) (inet, error) {
	if "" == network {
		return nil, fmt.Errorf("Network type is invalid!")
	}

	maker := makers[network]
	if nil == maker {
		return nil, fmt.Errorf("Ther is no valid maker for %s type", network)
	}

	return maker(network), nil
}

var makers map[string]netmaker = map[string]netmaker{
	"tcp": func(network string) inet {
		return new(netTcp)
	},
}

type netTcp struct {
}

func (n *netTcp) Listen(network string, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func (n *netTcp) Dial(network string, address string) (net.Conn, error) {
	return net.Dial(network, address)
}
