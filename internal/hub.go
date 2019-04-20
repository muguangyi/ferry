// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/muguangyi/gounite/network"
)

type HubSink struct {
}

func (hub *HubSink) OnConnected(p network.IPeer) {

}

func (hub *HubSink) OnClosed(p network.IPeer) {

}

func (hub *HubSink) OnPacket(p network.IPeer, obj interface{}) {

}
