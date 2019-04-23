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
	r.result = make(chan interface{})
	r.err = make(chan error)

	return r
}

type rpc struct {
	index  int64
	union  *Union
	result chan interface{}
	err    chan error
}

func (r *rpc) call(id string, name string, args ...interface{}) error {
	p := r.union.remoteUnits[id]
	if nil != p {
		req := &jsonPack{
			id: RPC_REQUEST,
			p: &protoRpcRequest{
				index:      r.index,
				unitName:   id,
				method:     name,
				args:       args,
				withResult: false,
			},
		}
		p.Send(req)
	}

	return nil
}

func (r *rpc) callWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	p := r.union.remoteUnits[id]
	if nil != p {
		req := &jsonPack{
			id: RPC_REQUEST,
			p: &protoRpcRequest{
				index:      r.index,
				unitName:   id,
				method:     name,
				args:       args,
				withResult: false,
			},
		}
		p.Send(req)

		result := <-r.result
		err := <-r.err
		close(r.result)
		close(r.err)
		return result, err
	}

	return nil, nil
}

func (r *rpc) callback(result interface{}, err error) {
	r.result <- result
	r.err <- err
}
