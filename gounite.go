// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gounite

import (
	"github.com/muguangyi/gounite/internal"
	"github.com/muguangyi/gounite/unit"
)

func Run(hubAddr string, units ...unit.IUnit) {
	union := &internal.Union{}
	union.Run(hubAddr, units)
}

func RunHub(hubAddr string) {
	hub := &internal.Hub{}
	hub.Run(hubAddr)
}
