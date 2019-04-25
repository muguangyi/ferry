// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unite

import (
	"sync"

	"github.com/muguangyi/unite/chancall"
)

func newUnit(control IUnitControl, discoverable bool) IUnit {
	u := new(unit)
	u.control = control
	u.discoverable = discoverable
	u.depends = make([]string, 0)
	u.callee = chancall.NewCallee(control)
	u.closeSig = make(chan bool, 1)

	return u
}

type unit struct {
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
	return u.union.call(id, name, args...)
}

func (u *unit) CallWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	return u.union.callWithResult(id, name, args...)
}

func (u *unit) SetTimeout(name string, timeout float32) {
	u.callee.SetTimeout(name, timeout)
}

func run(u *unit) {
	// u.control.OnUpdate(u.closeSig)
	u.wg.Done()
}

const (
	DEFAULT_TIMEOUT float32 = 1.0
)
