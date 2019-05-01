// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

import (
	"fmt"
	"log"
)

// ISignal interface.
type ISignal interface {
	// Setup signal booking.
	OnInit(signaler ISignaler)

	// Could start signal logic, like RPC etc.
	OnStart()

	// Clean up all.
	OnDestroy()
}

// ISignaler interface.
type ISignaler interface {
	// Setup dependency for other singal.
	Book(name string)

	// Get imported signal visitor.
	Visit(name string) interface{}

	// Call method with args, and no return value.
	Call(name string, method string, args ...interface{}) error

	// Call method with args, and has return values.
	CallWithResult(name string, method string, args ...interface{}) ([]interface{}, error)

	// Set target method with timeout duration.
	SetTimeout(method string, timeout float32)
}

// Startup run a signaler with target hub addr, customize union name for tracking, and
// all signals running in this signaler.
func Startup(hubAddr string, unionName string, signalers ...ISignaler) {
	union := newUnion(unionName, signalers...)
	union.run(hubAddr)
	wait(union)
}

// Serve run a hub with addr, black list for ports to avoid allocing to unions.
func Serve(hubAddr string, blackPorts ...int) {
	hub := newHub()
	hub.run(hubAddr, blackPorts...)
	wait(hub)
}

// Close all containers including hub or union
func Close() {
	destroy()
}

// NewSignaler create ISignaler with kernel signal which should implement ISignal interface.
func NewSignaler(id string, signal interface{}, discoverable bool) ISignaler {
	return newSignaler(id, signal, discoverable)
}

// Register signal id with proxy maker func.
func Register(id string, maker interface{}) bool {
	register(id, maker)
	return true
}

// Signal is base struct for all ISignal to compose
type Signal struct {
	signaler ISignaler
}

// OnInit initialize signal for other dependencies.
func (s *Signal) OnInit(signaler ISignaler) {
	s.signaler = signaler
}

// OnStart start signal logic, etc.
func (s *Signal) OnStart() {
}

// OnDestroy clean up all staff.
func (s *Signal) OnDestroy() {
}

// Book a signal with it's name.
func (s *Signal) Book(name string) {
	if nil != s.signaler {
		s.signaler.Book(name)
	} else {
		log.Fatal("ISignal not initialized, please make sure OnInit is called!")
	}
}

// Visit a signal with it's name.
func (s *Signal) Visit(name string) interface{} {
	if nil != s.signaler {
		return s.signaler.Visit(name)
	}

	log.Fatal("ISignal not initialized, please make sure OnInit is called!")
	return nil
}

// Call call target signal method with args by it's name, but without return value.
func (s *Signal) Call(name string, method string, args ...interface{}) error {
	if nil != s.signaler {
		return s.signaler.Call(name, method, args...)
	}

	log.Fatal("ISignal not initialized, please make sure OnInit is called!")

	return fmt.Errorf("ISignal not initialized, please make sure OnInit is called!")
}

// CallWithResult call target signal method with args by it's name, and with return values.
func (s *Signal) CallWithResult(name string, method string, args ...interface{}) ([]interface{}, error) {
	if nil != s.signaler {
		return s.signaler.CallWithResult(name, method, args...)
	}

	log.Fatal("ISignal not initialized, please make sure OnInit is called!")

	return nil, fmt.Errorf("ISignal not initialized, please make sure OnInit is called!")
}

// SetTimeout set method timeout value.
func (s *Signal) SetTimeout(method string, timeout float32) {
	if nil != s.signaler {
		s.signaler.SetTimeout(method, timeout)
	} else {
		log.Fatal("ISignal not initialized, please make sure OnInit is called!")
	}
}
