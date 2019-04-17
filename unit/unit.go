// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unit

import (
	"sync"
)

type unit struct {
	iu       IUnit
	closeSig chan bool
	wg       sync.WaitGroup
}

var units []*unit

func Load(iu IUnit) {
	u := new(unit)
	u.iu = iu
	u.closeSig = make(chan bool, 1)

	units = append(units, u)
}

func Init() {
	for i := 0; i < len(units); i++ {
		units[i].iu.OnInit()
	}

	for i := 0; i < len(units); i++ {
		u := units[i]
		u.wg.Add(1)
		go run(u)
	}
}

func run(u *unit) {
	u.iu.OnUpdate(u.closeSig)
	u.wg.Done()
}
