// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"github.com/muguangyi/gounite/framework"
)

func Run(hubAddr string, unionName string, units ...framework.IUnit) {
	union := framework.NewUnion(unionName, units...)
	union.Run(hubAddr)
}

func RunHub(hubAddr string) {
	hub := framework.NewHub()
	hub.Run(hubAddr)
}
