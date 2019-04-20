// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/muguangyi/gounite/network"
)

type ContainerSink struct {
}

func (c *ContainerSink) OnConnected(p network.IPeer) {
	if p.IsSelf() {
		// TODO:
		// 1.Register units to hub
	}
}

func (c *ContainerSink) OnClosed(p network.IPeer) {

}

func (c *ContainerSink) OnPacket(p network.IPeer, obj interface{}) {

}
