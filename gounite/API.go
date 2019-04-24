// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

// IUnitControl
type IUnitControl interface {
	// OnInit, setup method binding or unit importing
	OnInit(u IUnit)

	// OnStart, could start control logic, like RPC etc.
	OnStart()

	// OnDestroy, clean up all
	OnDestroy()
}

// IUnit
type IUnit interface {
	// Import, setup dependency for other units
	Import(id string)

	// Call, call method with args, and no return value
	Call(id string, name string, args ...interface{}) error

	// CallWithResult, call method with args, and has return value
	CallWithResult(id string, name string, args ...interface{}) (interface{}, error)

	// BindCall, bind method name with handle function
	BindCall(name string, function interface{})

	// BindCallWithTimeout, bind method with handle function, also with timeout
	BindCallWithTimeout(name string, function interface{}, timeout float32)
}

// Run, run an union with target hub addr, customize union name for tracking, and
// all units running in this union
func Run(hubAddr string, unionName string, units ...IUnit) {
	union := NewUnion(unionName, units...)
	union.Run(hubAddr)
}

// RunHub, run a hub with addr, black list for ports to avoid allocing to unions
func RunHub(hubAddr string, blackPorts ...int) {
	hub := NewHub()
	hub.Run(hubAddr, blackPorts...)
}

// NewUnit, new IUnit with IUnitControl object
func NewUnit(id string, control IUnitControl, discoverable bool) IUnit {
	return newUnit(id, control, discoverable)
}
