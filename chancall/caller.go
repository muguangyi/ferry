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

func (c *caller) Call(method string, args ...interface{}) error {
	err := c.call(&callRequest{
		method:       method,
		args:         args,
		callResponse: c.callResponse,
		done:         false,
	}, true)
	if nil != err {
		return err
	}

	response := <-c.callResponse
	return response.err
}

func (c *caller) CallWithResult(method string, args ...interface{}) ([]interface{}, error) {
	err := c.call(&callRequest{
		method:       method,
		args:         args,
		callResponse: c.callResponse,
		done:         false,
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
