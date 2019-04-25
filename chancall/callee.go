// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

import (
	"fmt"
	"time"
)

type callee struct {
	meta        *meta
	callRequest chan *callRequest
	functions   map[string]*fcall
}

func (c *callee) Name() string {
	return c.meta.name
}

func (c *callee) SetTimeout(name string, timeout float32) {
	c.meta.setTimeout(name, timeout)
}

func (c *callee) handling() {
	for {
		err := c.process(<-c.callRequest)
		if nil != err {
			panic(fmt.Sprintf("Invoke error %s", err.Error()))
		}
	}
}

func (c *callee) process(request *callRequest) (err error) {
	track(request, c.meta.timeout(request.method))
	result := c.meta.call(request.method, request.args...)
	return c.result(request, &callResponse{result: result})
}

func (c *callee) result(request *callRequest, response *callResponse) (err error) {
	if nil == request.callResponse {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	request.Lock()
	{
		if !request.done {
			request.done = true
			request.callResponse <- response
		}
	}
	request.Unlock()
	return
}

func track(request *callRequest, timeout float32) {
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		request.Lock()
		{
			if !request.done {
				request.done = true
				request.callResponse <- &callResponse{
					result: nil,
					err:    fmt.Errorf("[%s] function call timeout!", request.method),
				}
			}
		}
		request.Unlock()
	}()
}
