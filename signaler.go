// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/muguangyi/seek/chancall"
)

func newSignaler(name string, signal interface{}, discoverable bool) ISignaler {
	_, ok := signal.(ISignal)
	if !ok {
		panic(fmt.Sprintf("signal [%s] DOESNOT implement ISignal interface!", reflect.TypeOf(signal).Elem().Name()))
	}

	s := new(signaler)
	s.signal = signal.(ISignal)
	s.discoverable = discoverable
	s.depends = make([]string, 0)
	s.callee = chancall.NewCallee(name, signal)
	s.visiters = make(map[string]interface{})
	s.closeSig = make(chan bool, 1)

	return s
}

type signaler struct {
	signal       ISignal
	discoverable bool
	depends      []string
	callee       chancall.ICallee
	union        *union
	visiters     map[string]interface{}
	closeSig     chan bool
	wg           sync.WaitGroup
}

func (s *signaler) Book(name string) {
	if nil == s.union.localSignalers[name] {
		s.depends = append(s.depends, name)
	}
}

func (s *signaler) Visit(name string) interface{} {
	visitor := s.visiters[name]
	if nil == visitor {
		var ok bool
		if visitor, ok = tryMake(name, s); ok {
			s.visiters[name] = visitor
		}
	}

	return visitor
}

func (s *signaler) Call(name string, method string, args ...interface{}) error {
	return s.union.call(name, method, args...)
}

func (s *signaler) CallWithResult(name string, method string, args ...interface{}) ([]interface{}, error) {
	return s.union.callWithResult(name, method, args...)
}

func (s *signaler) SetTimeout(method string, timeout float32) {
	s.callee.SetTimeout(method, timeout)
}

func run(s *signaler) {
	// u.control.OnUpdate(u.closeSig)
	s.wg.Done()
}

const (
	cDefaultTimeout float32 = 1.0
)
