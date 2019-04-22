// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

type IUnitControl interface {
	OnInit(u IUnit)
	OnStart()
	OnDestroy()
	OnUpdate(closeSig chan bool)
}

type IUnit interface {
	Import(id string)
	Call(id string, name string, args ...interface{}) error
	CallWithResult(id string, name string, args ...interface{}) (interface{}, error)
	BindCall(name string, function interface{})
}

func NewUnit(id string, control IUnitControl, discoverable bool) IUnit {
	return newUnit(id, control, discoverable)
}
