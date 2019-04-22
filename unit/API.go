// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unit

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

func Init() {
	for _, u := range units {
		u.control.OnInit(u)
	}

	// for _, u := range units {
	// 	u.wg.Add(1)
	// 	go run(u)
	// }
}

func Collect() []string {
	ids := make([]string, 0)
	for id := range units {
		ids = append(ids, id)
	}

	return ids
}

func Depends() []string {
	ids := make([]string, 0)
	for id := range depends {
		ids = append(ids, id)
	}

	return ids
}
