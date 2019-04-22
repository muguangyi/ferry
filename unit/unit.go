// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unit

import (
	"sync"

	"github.com/muguangyi/gounite/network"
	"github.com/muguangyi/gounite/rpc"
)

func newUnit(id string, control IUnitControl, discoverable bool) IUnit {
	u := new(unit)
	u.id = id
	u.control = control
	u.discoverable = discoverable
	u.callee = rpc.NewCallee()
	u.closeSig = make(chan bool, 1)

	units[id] = u

	return u
}

type unit struct {
	id           string
	control      IUnitControl
	discoverable bool
	callee       rpc.ICallee
	closeSig     chan bool
	wg           sync.WaitGroup
}

type depend struct {
	addr string
	peer network.IPeer
}

func (u *unit) Import(id string) {
	if nil == depends[id] && nil == units[id] {
		depends[id] = new(depend)
	}
}

func (u *unit) Call(id string, name string, args ...interface{}) error {
	target := units[id]
	if nil != target {
		return rpc.NewCaller(target.callee).Call(name, args)
	} else {
		// TODO:
		// 1. Find unit in unions
		return nil
	}
}

func (u *unit) CallWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	target := units[id]
	if nil != target {
		return rpc.NewCaller(target.callee).CallWithResult(name, args)
	} else {
		// TODO:
		// 1. Find unit in unions
		return nil, nil
	}
}

func (u *unit) BindCall(name string, function interface{}) {
	u.callee.Bind(name, function)
}

var units map[string]*unit = make(map[string]*unit)
var depends map[string]*depend = make(map[string]*depend)

func run(u *unit) {
	u.control.OnUpdate(u.closeSig)
	u.wg.Done()
}
