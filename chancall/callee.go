// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

import (
	"fmt"
)

type callee struct {
	callRequest chan *callRequest
	functions   map[string]interface{}
}

func (c *callee) Bind(name string, function interface{}) {
	if _, ok := c.functions[name]; ok {
		panic(fmt.Sprintf("Function %v has already been binded!", name))
	}

	c.functions[name] = function
}

func (c *callee) Handling() {
	for {
		err := c.process(<-c.callRequest)
		if nil != err {
			panic(fmt.Sprintf("Invoke error %s", err.Error()))
		}
	}
}

func (c *callee) process(request *callRequest) (err error) {
	switch request.function.(type) {
	case func([]interface{}):
		request.function.(func([]interface{}))(request.args)
		return c.result(request, &callResponse{})
	case func([]interface{}) interface{}:
		result := request.function.(func([]interface{}) interface{})(request.args)
		return c.result(request, &callResponse{result: result})
	case func([]interface{}) []interface{}:
		result := request.function.(func([]interface{}) []interface{})(request.args)
		return c.result(request, &callResponse{result: result})
	}

	panic("bug")
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

	request.callResponse <- response
	return
}

func (c *callee) getFunction(name string, methodType int) (function interface{}, err error) {
	function = c.functions[name]
	if nil == function {
		err = fmt.Errorf("Function %v is not binded!", name)
		return
	}

	var ok bool
	switch methodType {
	case 0:
		_, ok = function.(func([]interface{}))
	case 1:
		_, ok = function.(func([]interface{}) interface{})
	case 2:
		_, ok = function.(func([]interface{}) []interface{})
	}

	if !ok {
		err = fmt.Errorf("Function %v type mismatched!", methodType)
	}

	return
}
