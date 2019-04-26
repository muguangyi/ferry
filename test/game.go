// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/muguangyi/unite/unite"
)

type IGame interface {
	Start(level string)
}

func newGame(wg *sync.WaitGroup) IGame {
	c := new(game)
	c.wg = wg

	return c
}

type game struct {
	unite.UnitControl
	wg *sync.WaitGroup
}

func (g *game) OnInit(u unite.IUnit) {
	g.UnitControl.OnInit(u)
	g.Import("IMath")
}

func (g *game) OnStart() {
	err := g.Call("IMath", "Print", "Hello World!")
	if nil != err {
		fmt.Println("error:", err.Error())
	}

	result, err := g.CallWithResult("IMath", "Add", 1, 2)
	if nil != err {
		fmt.Println("error:", err.Error())
	} else {
		fmt.Println("add result:", result)
	}

	g.wg.Done()
}

func (g *game) Start(level string) {
	fmt.Println("start...")
	fmt.Println("game started:", level)
}
