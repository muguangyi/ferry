// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"github.com/muguangyi/gounite/framework"
)

// Run, run an union with target hub addr, customize union name for tracking, and
// all units running in this union
func Run(hubAddr string, unionName string, units ...framework.IUnit) {
	union := framework.NewUnion(unionName, units...)
	union.Run(hubAddr)
}

// RunHub, run a hub with addr, black list for ports to avoid allocing to unions
func RunHub(hubAddr string, blackPorts ...int) {
	hub := framework.NewHub()
	hub.Run(hubAddr, blackPorts...)
}
