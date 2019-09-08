// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

// IFeature interface.
type IFeature interface {
	// Setup feature booking. Return what features it depends.
	OnInit() []string

	// Could start feature logic, like RPC etc.
	OnStart(s ISlot)

	// Clean up all.
	OnDestroy(s ISlot)
}

// ISlot interface.
type ISlot interface {
	// Get imported feature visitor.
	Visit(name string) interface{}

	// Call method with args, and no return value.
	Call(name string, method string, args ...interface{}) error

	// Call method with args, and has return values.
	CallWithResult(name string, method string, args ...interface{}) ([]interface{}, error)

	// Set target method with timeout duration.
	SetTimeout(method string, timeout float32)
}

// Startup run a dock with target hub addr, customize dock name for tracking, and
// all features running in this dock.
func Startup(hubAddr string, dockName string, slots ...ISlot) {
	dock := newDock(dockName, slots...)
	dock.run(hubAddr)
	wait(dock)
}

// Serve run a hub with addr, black list for ports to avoid allocing to docks.
func Serve(hubAddr string, blackPorts ...int) {
	hub := newHub()
	hub.run(hubAddr, blackPorts...)
	wait(hub)
}

// Close all containers including hub or dock.
func Close() {
	destroy()
}

// Carry an ISlot object with kernel feature that should implement IFeature interface.
func Carry(id string, feature interface{}, discoverable bool) ISlot {
	return newSlot(id, feature, discoverable)
}

// Register feature id with proxy maker func.
func Register(id string, maker interface{}) bool {
	register(id, maker)
	return true
}

// Feature is base struct for all IFeature to compose
type Feature struct {
}

// OnInit initialize feature for other dependencies.
func (f *Feature) OnInit() []string {
	return nil
}

// OnStart start feature logic, etc.
func (f *Feature) OnStart(s ISlot) {
}

// OnDestroy clean up all staff.
func (f *Feature) OnDestroy(s ISlot) {
}
