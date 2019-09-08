// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"fmt"
	"time"
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

func (r *rpc) call(dock *dock, name string, method string, args ...interface{}) error {
	peer := dock.queryRemoteSlot(name)
	if peer != nil {
		req := &packer{
			Id: cRpcRequest,
			P: &protoRpcRequest{
				Index:      r.index,
				SlotId:     name,
				Method:     method,
				Args:       args,
				WithResult: false,
			},
		}
		peer.Send(req)
	} else {
		// TODO: Handle defer rpc call.
		r.callback(&ret{result: nil, err: fmt.Errorf("NO [%s] slot exist!", name)})
	}

	ret := <-r.ret
	close(r.ret)

	return ret.err
}

func (r *rpc) callWithResult(dock *dock, name string, method string, args ...interface{}) ([]interface{}, error) {
	peer := dock.queryRemoteSlot(name)
	if peer != nil {
		req := &packer{
			Id: cRpcRequest,
			P: &protoRpcRequest{
				Index:      r.index,
				SlotId:     name,
				Method:     method,
				Args:       args,
				WithResult: true,
			},
		}
		peer.Send(req)
	} else {
		// TODO: Handle defer rpc call.
		r.callback(&ret{result: nil, err: fmt.Errorf("NO [%s] slot exist!", name)})
	}

	ret := <-r.ret
	close(r.ret)

	return ret.result, ret.err
}

func (r *rpc) callback(ret *ret) {
	r.ret <- ret
}
