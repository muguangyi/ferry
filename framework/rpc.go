// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

import (
	"time"
)

func newRpc(union *Union) *rpc {
	r := new(rpc)
	r.index = time.Now().UnixNano()
	r.union = union
	r.result = make(chan interface{}, 1)

	return r
}

type rpc struct {
	index  int64
	union  *Union
	result chan interface{}
}

func (r *rpc) call(id string, name string, args ...interface{}) error {
	p := r.union.remoteUnits[id]
	if nil != p {
		req := &jsonPack{
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
	}

	return nil
}

func (r *rpc) callWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	p := r.union.remoteUnits[id]
	if nil != p {
		r.union.invoke(r)
		req := &jsonPack{
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

		result := <-r.result
		close(r.result)
		return result, nil
	}

	return nil, nil
}

func (r *rpc) callback(result interface{}) {
	r.result <- result
}
