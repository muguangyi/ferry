// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unit

import (
	"sync"
)

type unit struct {
	id           string
	control      IUnitControl
	discoverable bool
	closeSig     chan bool
	wg           sync.WaitGroup
}

var units []*unit

func run(u *unit) {
	u.control.OnUpdate(u.closeSig)
	u.wg.Done()
}
