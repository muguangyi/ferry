// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

// ICallee
type ICallee interface {
	// Bind, bind method name with corresponding function and timeout
	Bind(name string, function interface{}, timeout float32)
}

// ICaller
type ICaller interface {
	// Call, call "name" method followed args and no return value
	Call(name string, args ...interface{}) error

	// CallWithResult, call "name" method followed args and has one return value
	CallWithResult(name string, args ...interface{}) (interface{}, error)
}

// NewCallee, create a new callee
func NewCallee() ICallee {
	c := new(callee)
	c.callRequest = make(chan *callRequest, 2)
	c.functions = make(map[string]*fcall)
	go c.handling()

	return c
}

// NewCaller, create a caller for target callee
func NewCaller(c ICallee) ICaller {
	caller := new(caller)
	caller.callee = c.(*callee)
	caller.callResponse = make(chan *callResponse, 1)

	return caller
}
