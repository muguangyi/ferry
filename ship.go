// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ship

import (
	"fmt"
	"log"
)

// IFeature interface.
type IFeature interface {
	// Setup feature booking.
	OnInit(sandbox ISandbox)

	// Could start feature logic, like RPC etc.
	OnStart()

	// Clean up all.
	OnDestroy()
}

// ISandbox interface.
type ISandbox interface {
	// Setup dependency for other feature.
	Book(name string)

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
func Startup(hubAddr string, dockName string, sandboxes ...ISandbox) {
	dock := newDock(dockName, sandboxes...)
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

// Carry an ISandbox object with kernel feature that should implement IFeature interface.
func Carry(id string, feature interface{}, discoverable bool) ISandbox {
	return newSandbox(id, feature, discoverable)
}

// Register feature id with proxy maker func.
func Register(id string, maker interface{}) bool {
	register(id, maker)
	return true
}

// Feature is base struct for all IFeature to compose
type Feature struct {
	sandbox ISandbox
}

// OnInit initialize feature for other dependencies.
func (f *Feature) OnInit(sandbox ISandbox) {
	f.sandbox = sandbox
}

// OnStart start feature logic, etc.
func (f *Feature) OnStart() {
}

// OnDestroy clean up all staff.
func (f *Feature) OnDestroy() {
}

// Book a feature with it's name.
func (f *Feature) Book(name string) {
	if nil != f.sandbox {
		f.sandbox.Book(name)
	} else {
		log.Fatal("IFeature not initialized, please make sure OnInit is called!")
	}
}

// Visit feature with it's name.
func (f *Feature) Visit(name string) interface{} {
	if nil != f.sandbox {
		return f.sandbox.Visit(name)
	}

	log.Fatal("IFeature not initialized, please make sure OnInit is called!")
	return nil
}

// Call call target feature method with args by it's name, but without return value.
func (f *Feature) Call(name string, method string, args ...interface{}) error {
	if nil != f.sandbox {
		return f.sandbox.Call(name, method, args...)
	}

	log.Fatal("IFeature not initialized, please make sure OnInit is called!")

	return fmt.Errorf("IFeature not initialized, please make sure OnInit is called!")
}

// CallWithResult call target feature method with args by it's name, and with return values.
func (f *Feature) CallWithResult(name string, method string, args ...interface{}) ([]interface{}, error) {
	if nil != f.sandbox {
		return f.sandbox.CallWithResult(name, method, args...)
	}

	log.Fatal("IFeature not initialized, please make sure OnInit is called!")

	return nil, fmt.Errorf("ISignal not initialized, please make sure OnInit is called!")
}

// SetTimeout set method timeout value.
func (f *Feature) SetTimeout(method string, timeout float32) {
	if nil != f.sandbox {
		f.sandbox.SetTimeout(method, timeout)
	} else {
		log.Fatal("IFeature not initialized, please make sure OnInit is called!")
	}
}
