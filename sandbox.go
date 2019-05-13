// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/muguangyi/ferry/chancall"
)

func newSandbox(name string, feature interface{}, discoverable bool) ISandbox {
	_, ok := feature.(IFeature)
	if !ok {
		panic(fmt.Sprintf("Feature [%s] DOES NOT implement IFeature interface!", reflect.TypeOf(feature).Elem().Name()))
	}

	s := new(sandbox)
	s.feature = feature.(IFeature)
	s.discoverable = discoverable
	s.depends = make([]string, 0)
	s.callee = chancall.NewCallee(name, feature)
	s.visiters = make(map[string]interface{})
	s.closeSig = make(chan bool, 1)

	return s
}

type sandbox struct {
	feature      IFeature
	discoverable bool
	depends      []string
	callee       chancall.ICallee
	dock         *dock
	visiters     map[string]interface{}
	closeSig     chan bool
	wg           sync.WaitGroup
}

func (s *sandbox) Book(name string) {
	if nil == s.dock.sandboxes[name] {
		s.depends = append(s.depends, name)
	}
}

func (s *sandbox) Visit(name string) interface{} {
	visitor := s.visiters[name]
	if nil == visitor {
		var ok bool
		if visitor, ok = tryMake(name, s); ok {
			s.visiters[name] = visitor
		}
	}

	return visitor
}

func (s *sandbox) Call(name string, method string, args ...interface{}) error {
	return s.dock.call(name, method, args...)
}

func (s *sandbox) CallWithResult(name string, method string, args ...interface{}) ([]interface{}, error) {
	return s.dock.callWithResult(name, method, args...)
}

func (s *sandbox) SetTimeout(method string, timeout float32) {
	s.callee.SetTimeout(method, timeout)
}

func run(s *sandbox) {
	// u.control.OnUpdate(u.closeSig)
	s.wg.Done()
}

const (
	cDefaultTimeout float32 = 1.0
)
