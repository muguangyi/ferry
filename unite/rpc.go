// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unite

import (
	"time"

	"github.com/muguangyi/unite/network"
)

func newRpc() *rpc {
	r := new(rpc)
	r.index = time.Now().UnixNano()
	r.ret = make(chan *ret, 1)

	return r
}

type rpc struct {
	index int64
	ret   chan *ret
}

type ret struct {
	result []interface{}
	err    error
}

func (r *rpc) call(peer network.IPeer, id string, name string, args ...interface{}) error {
	req := &packer{
		Id: RPC_REQUEST,
		P: &protoRpcRequest{
			Index:      r.index,
			UnitId:     id,
			Method:     name,
			Args:       args,
			WithResult: false,
		},
	}
	peer.Send(req)

	ret := <-r.ret
	close(r.ret)

	return ret.err
}

func (r *rpc) callWithResult(peer network.IPeer, id string, name string, args ...interface{}) ([]interface{}, error) {
	req := &packer{
		Id: RPC_REQUEST,
		P: &protoRpcRequest{
			Index:      r.index,
			UnitId:     id,
			Method:     name,
			Args:       args,
			WithResult: true,
		},
	}
	peer.Send(req)

	ret := <-r.ret
	close(r.ret)

	return ret.result, ret.err
}

func (r *rpc) callback(ret *ret) {
	r.ret <- ret
}
