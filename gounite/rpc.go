// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"time"
)

func newRpc(union *union) *rpc {
	r := new(rpc)
	r.index = time.Now().UnixNano()
	r.union = union
	r.ret = make(chan *ret, 1)

	return r
}

type rpc struct {
	index int64
	union *union
	ret   chan *ret
}

type ret struct {
	result interface{}
	err    error
}

func (r *rpc) call(id string, name string, args ...interface{}) error {
	p := r.union.remoteUnits[id]
	if nil != p {
		r.union.invoke(r)
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
		p.Send(req)

		ret := <-r.ret
		close(r.ret)

		return ret.err
	}

	return nil
}

func (r *rpc) callWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	p := r.union.remoteUnits[id]
	if nil != p {
		r.union.invoke(r)
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
		p.Send(req)

		ret := <-r.ret
		close(r.ret)

		return ret.result, ret.err
	}

	return nil, nil
}

func (r *rpc) callback(ret *ret) {
	r.ret <- ret
}
