// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

// ICallee interface.
type ICallee interface {
	// Return callee's name.
	Name() string

	// Set target method timeout duration.
	SetTimeout(name string, timeout float32)
}

// ICaller interface.
type ICaller interface {
	// Call "name" method followed args and no return value.
	Call(name string, args ...interface{}) error

	// Call "name" method followed args and has return values.
	CallWithResult(name string, args ...interface{}) ([]interface{}, error)
}

// Create a new callee with unique name and target object.
func NewCallee(name string, target interface{}) ICallee {
	c := new(callee)
	c.meta = newMeta(name, target)
	c.callRequest = make(chan *callRequest, 2)
	c.functions = make(map[string]*fcall)
	go c.handling()

	return c
}

// Create a caller for target callee.
func NewCaller(c ICallee) ICaller {
	caller := new(caller)
	caller.callee = c.(*callee)
	caller.callResponse = make(chan *callResponse, 1)

	return caller
}
