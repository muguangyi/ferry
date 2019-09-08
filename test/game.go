// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/muguangyi/ferry"
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
	ferry.Feature
	wg *sync.WaitGroup
}

func (g *game) OnInit() []string {
	return []string{"IMath"}
}

func (g *game) OnStart(s ferry.ISlot) {
	math := s.Visit("IMath").(IMath)
	math.Print("Hello World!")

	result := math.Add(1, 2)
	fmt.Println("add result:", result)

	g.wg.Done()
}

func (g *game) Start(level string) {
	fmt.Println("start...")
	fmt.Println("game started:", level)
}
