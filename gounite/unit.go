// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"sync"

	"github.com/muguangyi/gounite/chancall"
)

func newUnit(id string, control IUnitControl, discoverable bool) IUnit {
	u := new(unit)
	u.id = id
	u.control = control
	u.discoverable = discoverable
	u.depends = make([]string, 0)
	u.callee = chancall.NewCallee()
	u.closeSig = make(chan bool, 1)

	return u
}

type unit struct {
	id           string
	control      IUnitControl
	discoverable bool
	depends      []string
	callee       chancall.ICallee
	union        *union
	closeSig     chan bool
	wg           sync.WaitGroup
}

func (u *unit) Import(id string) {
	if nil == u.union.localUnits[id] {
		u.depends = append(u.depends, id)
	}
}

func (u *unit) Call(id string, name string, args ...interface{}) error {
	target := u.union.localUnits[id]
	if nil != target {
		return chancall.NewCaller(target.callee).Call(name, args...)
	} else {
		rpc := newRpc(u.union)
		return rpc.call(id, name, args...)
	}
}

func (u *unit) CallWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	target := u.union.localUnits[id]
	if nil != target {
		return chancall.NewCaller(target.callee).CallWithResult(name, args...)
	} else {
		rpc := newRpc(u.union)
		return rpc.callWithResult(id, name, args...)
	}
}

func (u *unit) BindCall(name string, function interface{}) {
	u.callee.Bind(name, function, DEFAULT_TIMEOUT)
}

func (u *unit) BindCallWithTimeout(name string, function interface{}, timeout float32) {
	u.callee.Bind(name, function, timeout)
}

func run(u *unit) {
	// u.control.OnUpdate(u.closeSig)
	u.wg.Done()
}

const (
	DEFAULT_TIMEOUT float32 = 1.0
)
