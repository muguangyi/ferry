// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

import (
	"fmt"
	"time"
)

type caller struct {
	callee       *callee
	callResponse chan *callResponse
}

func (c *caller) Call(name string, args ...interface{}) error {
	function, timeout, err := c.callee.search(name, 0)
	if nil != err {
		return err
	}

	req := &callRequest{
		function:     function,
		args:         args,
		callResponse: c.callResponse,
		done:         false,
	}
	track(name, req, timeout)

	err = c.call(req, true)
	if nil != err {
		return err
	}

	response := <-c.callResponse
	return response.err
}

func (c *caller) CallWithResult(name string, args ...interface{}) (interface{}, error) {
	function, timeout, err := c.callee.search(name, 1)
	if nil != err {
		return nil, err
	}

	req := &callRequest{
		function:     function,
		args:         args,
		callResponse: c.callResponse,
		done:         false,
	}
	track(name, req, timeout)

	err = c.call(req, true)
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

func track(name string, request *callRequest, timeout float32) {
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		request.Lock()
		{
			if !request.done {
				request.done = true
				request.callResponse <- &callResponse{
					result: nil,
					err:    fmt.Errorf("[%s] function call timeout!", name),
				}
			}
		}
		request.Unlock()
	}()
}
