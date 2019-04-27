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

func newSignaler(id string, signal interface{}, discoverable bool) ISignaler {
	_, ok := signal.(ISignal)
	if !ok {
		panic(fmt.Sprintf("signal [%s] DOESNOT implement ISignal interface!", reflect.TypeOf(signal).Elem().Name()))
	}

	s := new(signaler)
	s.signal = signal.(ISignal)
	s.discoverable = discoverable
	s.depends = make([]string, 0)
	s.callee = chancall.NewCallee(id, signal)
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

func (s *signaler) Book(id string) {
	if nil == s.union.localSignalers[id] {
		s.depends = append(s.depends, id)
	}
}

func (s *signaler) Visit(id string) interface{} {
	visitor := s.visiters[id]
	if nil == visitor {
		if visitor, ok := tryMake(id, s); ok {
			s.visiters[id] = visitor
		}
	}

	return visitor
}

func (s *signaler) Call(id string, name string, args ...interface{}) error {
	return s.union.call(id, name, args...)
}

func (s *signaler) CallWithResult(id string, name string, args ...interface{}) ([]interface{}, error) {
	return s.union.callWithResult(id, name, args...)
}

func (s *signaler) SetTimeout(name string, timeout float32) {
	s.callee.SetTimeout(name, timeout)
}

func run(s *signaler) {
	// u.control.OnUpdate(u.closeSig)
	s.wg.Done()
}

const (
	DEFAULT_TIMEOUT float32 = 1.0
)
