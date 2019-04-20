// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/muguangyi/gounite/network"
)

type UnionSink struct {
}

func (un *UnionSink) OnConnected(p network.IPeer) {
	if p.IsSelf() {
		// TODO:
		// 1.Register units to hub
	}
}

func (un *UnionSink) OnClosed(p network.IPeer) {

}

func (un *UnionSink) OnPacket(p network.IPeer, obj interface{}) {

}
