// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"time"
)

func newRpc() *rpc {
	return &rpc{index: time.Now().UnixNano(), req: nil, ret: make(chan *ret, 1)}
}

type rpc struct {
	index int64
	req   *protoRpcRequest
	ret   chan *ret
}

type ret struct {
	result []interface{}
	err    error
}

func (r *rpc) call(dock *dock, name string, method string, args ...interface{}) error {
	r.req = &protoRpcRequest{
		Index:      r.index,
		Slot:       name,
		Method:     method,
		Args:       args,
		WithResult: false,
	}

	dock.commit(r)

	ret := <-r.ret
	close(r.ret)

	return ret.err
}

func (r *rpc) callWithResult(dock *dock, name string, method string, args ...interface{}) ([]interface{}, error) {
	r.req = &protoRpcRequest{
		Index:      r.index,
		Slot:       name,
		Method:     method,
		Args:       args,
		WithResult: true,
	}

	dock.commit(r)

	ret := <-r.ret
	close(r.ret)

	return ret.result, ret.err
}

func (r *rpc) callback(ret *ret) {
	r.ret <- ret
}
