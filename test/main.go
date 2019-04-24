// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/muguangyi/gounite"
	"github.com/muguangyi/gounite/framework"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	gounite.RunHub("127.0.0.1:9999")

	gounite.Run("127.0.0.1:9999", "util",
		framework.NewUnit("math", &MathControl{}, true))

	gounite.Run("127.0.0.1:9999", "logic",
		framework.NewUnit("game", &GameControl{wg: &wg}, true),
		framework.NewUnit("lobby", &LobbyControl{wg: &wg}, true))

	wg.Wait()
	fmt.Println("Completed!")
}

type MathControl struct {
	unit framework.IUnit
}

func (math *MathControl) OnInit(u framework.IUnit) {
	math.unit = u
	u.BindCall("print", math.print)
	u.BindCall("add", math.add)
}

func (math *MathControl) OnStart() {

}

func (math *MathControl) OnDestroy() {

}

func (math *MathControl) print(args []interface{}) {
	fmt.Println("-----print")
	fmt.Println(args[0].(string))
}

func (math *MathControl) add(args []interface{}) interface{} {
	fmt.Println(fmt.Sprintf("-----add method called, %T, %T", args[0], args[1]))
	result := args[0].(float64) + args[1].(float64)
	return result
}

type LobbyControl struct {
	unit framework.IUnit
	wg   *sync.WaitGroup
}

func (l *LobbyControl) OnInit(u framework.IUnit) {
	l.unit = u
	u.Import("game")
}

func (l *LobbyControl) OnStart() {
	l.unit.Call("game", "start", "level1")
	l.wg.Done()
}

func (l *LobbyControl) OnDestroy() {

}

type GameControl struct {
	unit framework.IUnit
	wg   *sync.WaitGroup
}

func (g *GameControl) OnInit(u framework.IUnit) {
	g.unit = u
	u.Import("math")
	u.BindCall("start", g.start)
}

func (g *GameControl) OnStart() {
	g.unit.Call("math", "print", "Hello World!")
	result, _ := g.unit.CallWithResult("math", "add", 1, 2)
	fmt.Println("-----Math add result:", result)
	g.wg.Done()
}

func (g *GameControl) OnDestroy() {

}

func (g *GameControl) start(args []interface{}) {
	fmt.Println("--> game started:", args[0].(string))
}
