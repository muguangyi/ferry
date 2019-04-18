// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unit

type IUnitControl interface {
	OnInit()
	OnDestroy()
	OnUpdate(closeSig chan bool)
}

type IUnit interface {
}

func NewUnit(id string, control IUnitControl, discoverable bool) IUnit {
	u := new(unit)
	u.id = id
	u.control = control
	u.discoverable = discoverable
	u.closeSig = make(chan bool, 1)

	units = append(units, u)

	return u
}

func Init() {
	for i := 0; i < len(units); i++ {
		units[i].control.OnInit()
	}

	for i := 0; i < len(units); i++ {
		u := units[i]
		u.wg.Add(1)
		go run(u)
	}
}
