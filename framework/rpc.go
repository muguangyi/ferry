// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

import (
	"github.com/muguangyi/gounite/network"
)

var remoteUnits map[string]network.IPeer = make(map[string]network.IPeer)

func RemoteCall(id string, name string, args ...interface{}) error {
	p := remoteUnits[id]
	if nil != p {
		req := &jsonPack{
			id: RPC_REQUEST,
			p: &protoRpcRequest{
				unitName: id,
				method:   name,
				args:     args,
			},
		}
		p.Send(req)
	}

	return nil
}
