// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unite

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

	// SetTimeout, set target method with timeout duration
	SetTimeout(name string, timeout float32)
}

// Run, run an union with target hub addr, customize union name for tracking, and
// all units running in this union
func Run(hubAddr string, unionName string, units ...IUnit) {
	union := newUnion(unionName, units...)
	union.run(hubAddr)
}

// RunHub, run a hub with addr, black list for ports to avoid allocing to unions
func RunHub(hubAddr string, blackPorts ...int) {
	hub := newHub()
	hub.run(hubAddr, blackPorts...)
}

// NewUnit, new IUnit with IUnitControl object
func NewUnit(control IUnitControl, discoverable bool) IUnit {
	return newUnit(control, discoverable)
}

// UnitControl, base struct for all IUnitControl to compose
type UnitControl struct {
	unit IUnit
}

func (c *UnitControl) OnInit(u IUnit) {
	c.unit = u
}

func (c *UnitControl) OnStart() {
}

func (c *UnitControl) OnDestroy() {
}

func (c *UnitControl) Import(id string) {
	if nil != c.unit {
		c.unit.Import(id)
	} else {
		panic("IUnit not initialized, please make sure OnInit is called!")
	}
}

func (c *UnitControl) Call(id string, name string, args ...interface{}) error {
	if nil == c.unit {
		return c.unit.Call(id, name, args...)
	}

	panic("IUnit not initialized, please make sure OnInit is called!")
}

func (c *UnitControl) CallWithResult(id string, name string, args ...interface{}) (interface{}, error) {
	if nil != c.unit {
		return c.unit.CallWithResult(id, name, args...)
	}

	panic("IUnit not initialized, please make sure OnInit is called!")
}

func (c *UnitControl) SetTimeout(name string, timeout float32) {
	if nil != c.unit {
		c.unit.SetTimeout(name, timeout)
	} else {
		panic("IUnit not initialized, please make sure OnInit is called!")

	}
}
