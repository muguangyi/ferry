// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

import (
	"fmt"
)

type caller struct {
	callee       *callee
	callResponse chan *callResponse
}

func (c *caller) Call(name string, args ...interface{}) error {
	function, err := c.callee.getFunction(name, 0)
	if nil != err {
		return err
	}

	err = c.call(&callRequest{
		function:     function,
		args:         args,
		callResponse: c.callResponse,
	}, true)

	if nil != err {
		return err
	}

	response := <-c.callResponse
	return response.err
}

func (c *caller) CallWithResult(name string, args ...interface{}) (interface{}, error) {
	function, err := c.callee.getFunction(name, 1)
	if nil != err {
		return nil, err
	}

	err = c.call(&callRequest{
		function:     function,
		args:         args,
		callResponse: c.callResponse,
	}, true)

	if nil != err {
		return nil, err
	}

	response := <-c.callResponse
	return response.result, response.err
}

func (c *caller) call(request *callRequest, block bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if block {
		c.callee.callRequest <- request
	} else {
		select {
		case c.callee.callRequest <- request:
		default:
			err = fmt.Errorf("RPC channel full!")
		}
	}

	return
}
