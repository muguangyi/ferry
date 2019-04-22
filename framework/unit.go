// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

import (
	"sync"

	"github.com/muguangyi/gounite/rpc"
)

func newUnit(id string, control IUnitControl, discoverable bool) IUnit {
	u := new(unit)
	u.id = id
	u.control = control
	u.discoverable = discoverable
	u.callee = rpc.NewCallee()
	u.closeSig = make(chan bool, 1)

	localUnits[id] = u

	return u
}

type unit struct {
	id           string
	control      IUnitControl
	discoverable bool
	depends      []string
	callee       rpc.ICallee
	closeSig     chan bool
	wg           sync.WaitGroup
}

func (u *unit) Import(id string) {
	if nil == localUnits[id] {
		u.depends = append(u.depends, id)
	}
}

func (u *unit) Call(id string, name string, args ...interface{}) error {
	target := localUnits[id]
	if nil != target {
		return rpc.NewCaller(target.callee).Call(name, args)
	} else {
		return RemoteCall(id, name, args)
	}
}

func (u *unit) CallWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	target := localUnits[id]
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

var localUnits map[string]*unit = make(map[string]*unit)

func localInit() {
	for _, u := range localUnits {
		u.control.OnInit(u)
	}

	// for _, u := range units {
	// 	u.wg.Add(1)
	// 	go run(u)
	// }
}

func localStart() {
	for _, u := range localUnits {
		u.control.OnStart()
	}
}

func localCollect() []string {
	ids := make([]string, 0)
	for id := range remoteUnits {
		ids = append(ids, id)
	}

	return ids
}

func localDepends() []string {
	ids := make([]string, 0)
	for _, v := range localUnits {
		ids = append(ids, v.depends...)
	}

	return ids
}

func run(u *unit) {
	u.control.OnUpdate(u.closeSig)
	u.wg.Done()
}
