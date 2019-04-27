// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

import (
	"fmt"
	"log"
)

// ISignal
type ISignal interface {
	// OnInit, setup signal booking
	OnInit(signaler ISignaler)

	// OnStart, could start signal logic, like RPC etc.
	OnStart()

	// OnDestroy, clean up all
	OnDestroy()
}

// ISignaler
type ISignaler interface {
	// Book, setup dependency for other singal
	Book(id string)

	// Visit, get imported signal visitor
	Visit(id string) interface{}

	// Call, call method with args, and no return value
	Call(id string, name string, args ...interface{}) error

	// CallWithResult, call method with args, and has return value
	CallWithResult(id string, name string, args ...interface{}) ([]interface{}, error)

	// SetTimeout, set target method with timeout duration
	SetTimeout(name string, timeout float32)
}

// Startup, run a signaler with target hub addr, customize union name for tracking, and
// all signals running in this signaler
func Startup(hubAddr string, unionName string, signalers ...ISignaler) {
	union := newUnion(unionName, signalers...)
	union.run(hubAddr)
}

// RunHub, run a hub with addr, black list for ports to avoid allocing to unions
func Serve(hubAddr string, blackPorts ...int) {
	hub := newHub()
	hub.run(hubAddr, blackPorts...)
}

// NewSignaler, new ISignaler with kernel signal which should implement ISignal interface
func NewSignaler(id string, signal interface{}, discoverable bool) ISignaler {
	return newSignaler(id, signal, discoverable)
}

// Register, register signal id with proxy maker func
func Register(id string, maker interface{}) bool {
	register(id, maker)
	return true
}

// Signal, base struct for all ISignal to compose
type Signal struct {
	signaler ISignaler
}

func (s *Signal) OnInit(signaler ISignaler) {
	s.signaler = signaler
}

func (s *Signal) OnStart() {
}

func (s *Signal) OnDestroy() {
}

func (s *Signal) Book(id string) {
	if nil != s.signaler {
		s.signaler.Book(id)
	} else {
		log.Fatal("ISignal not initialized, please make sure OnInit is called!")
	}
}

func (s *Signal) Visit(id string) interface{} {
	if nil != s.signaler {
		return s.signaler.Visit(id)
	}

	log.Fatal("ISignal not initialized, please make sure OnInit is called!")
	return nil
}

func (s *Signal) Call(id string, name string, args ...interface{}) error {
	if nil != s.signaler {
		return s.signaler.Call(id, name, args...)
	}

	log.Fatal("ISignal not initialized, please make sure OnInit is called!")

	return fmt.Errorf("ISignal not initialized, please make sure OnInit is called!")
}

func (s *Signal) CallWithResult(id string, name string, args ...interface{}) ([]interface{}, error) {
	if nil != s.signaler {
		return s.signaler.CallWithResult(id, name, args...)
	}

	log.Fatal("ISignal not initialized, please make sure OnInit is called!")

	return nil, fmt.Errorf("ISignal not initialized, please make sure OnInit is called!")
}

func (s *Signal) SetTimeout(name string, timeout float32) {
	if nil != s.signaler {
		s.signaler.SetTimeout(name, timeout)
	} else {
		log.Fatal("ISignal not initialized, please make sure OnInit is called!")
	}
}
