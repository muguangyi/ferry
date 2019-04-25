// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/muguangyi/unite/unite"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	unite.RunHub("127.0.0.1:9999")

	unite.Run("127.0.0.1:9999", "util",
		unite.NewUnit(&MathControl{}, true))

	unite.Run("127.0.0.1:9999", "logic",
		unite.NewUnit(&GameControl{wg: &wg}, true),
		unite.NewUnit(&LobbyControl{wg: &wg}, true))

	wg.Wait()
	fmt.Println("Completed!")
}

type MathControl struct {
	unite.UnitControl
}

func (math *MathControl) Print(msg string) {
	time.Sleep(5 * time.Second)
	fmt.Println("print...")
	fmt.Println(msg)
}

func (math *MathControl) Add(x float64, y float64) interface{} {
	// time.Sleep(5 * time.Second)
	fmt.Println("add...")
	result := x + y
	return result
}

type LobbyControl struct {
	unite.UnitControl
	wg *sync.WaitGroup
}

func (l *LobbyControl) OnInit(u unite.IUnit) {
	l.UnitControl.OnInit(u)
	l.Import("GameControl")
}

func (l *LobbyControl) OnStart() {
	l.Call("GameControl", "Start", "level1")
	l.wg.Done()
}

type GameControl struct {
	unite.UnitControl
	wg *sync.WaitGroup
}

func (g *GameControl) OnInit(u unite.IUnit) {
	g.UnitControl.OnInit(u)
	g.Import("MathControl")
}

func (g *GameControl) OnStart() {
	err := g.Call("MathControl", "Print", "Hello World!")
	if nil != err {
		fmt.Println("error:", err.Error())
	}

	result, err := g.CallWithResult("MathControl", "Add", 1, 2)
	if nil != err {
		fmt.Println("error:", err.Error())
	} else {
		fmt.Println("add result:", result)
	}

	g.wg.Done()
}

func (g *GameControl) Start(level string) {
	fmt.Println("start...")
	fmt.Println("game started:", level)
}
