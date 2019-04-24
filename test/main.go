// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	// "time"

	"github.com/muguangyi/unite/unite"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	unite.RunHub("127.0.0.1:9999")

	unite.Run("127.0.0.1:9999", "util",
		unite.NewUnit("math", &MathControl{}, true))

	unite.Run("127.0.0.1:9999", "logic",
		unite.NewUnit("game", &GameControl{wg: &wg}, true),
		unite.NewUnit("lobby", &LobbyControl{wg: &wg}, true))

	wg.Wait()
	fmt.Println("Completed!")
}

type MathControl struct {
	unit unite.IUnit
}

func (math *MathControl) OnInit(u unite.IUnit) {
	math.unit = u
	u.BindCall("print", math.print)
	u.BindCall("add", math.add)
}

func (math *MathControl) OnStart() {

}

func (math *MathControl) OnDestroy() {

}

func (math *MathControl) print(args []interface{}) {
	// time.Sleep(5 * time.Second)
	fmt.Println("print...")
	fmt.Println(args[0].(string))
}

func (math *MathControl) add(args []interface{}) interface{} {
	// time.Sleep(5 * time.Second)
	fmt.Println("add...")
	result := args[0].(float64) + args[1].(float64)
	return result
}

type LobbyControl struct {
	unit unite.IUnit
	wg   *sync.WaitGroup
}

func (l *LobbyControl) OnInit(u unite.IUnit) {
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
	unit gounite.IUnit
	wg   *sync.WaitGroup
}

func (g *GameControl) OnInit(u unite.IUnit) {
	g.unit = u
	u.Import("math")
	u.BindCall("start", g.start)
}

func (g *GameControl) OnStart() {
	err := g.unit.Call("math", "print", "Hello World!")
	if nil != err {
		fmt.Println("error:", err.Error())
	}

	result, err := g.unit.CallWithResult("math", "add", 1, 2)
	if nil != err {
		fmt.Println("error:", err.Error())
	} else {
		fmt.Println("add result:", result)
	}

	g.wg.Done()
}

func (g *GameControl) OnDestroy() {

}

func (g *GameControl) start(args []interface{}) {
	fmt.Println("start...")
	fmt.Println("game started:", args[0].(string))
}
