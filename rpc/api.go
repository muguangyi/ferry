// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package rpc

type ICallee interface {
	Bind(name string, function interface{})
	Handling()
}

type ICaller interface {
	Call(name string, args ...interface{}) error
	CallWithResult(name string, args ...interface{}) (interface{}, error)
}

func NewCallee() ICallee {
	c := new(callee)
	c.callRequest = make(chan *callRequest, 1)
	c.functions = make(map[string]interface{})

	return c
}

func NewCaller(c ICallee) ICaller {
	caller := new(caller)
	caller.callee = c.(*callee)
	caller.callResponse = make(chan *callResponse, 1)

	return caller
}
