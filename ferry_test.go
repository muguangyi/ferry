// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry_test

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/muguangyi/ferry"
	"github.com/muguangyi/ferry/network"
)

type ILogger interface {
	Log(v interface{})
}

type logger struct {
	ferry.Feature
}

func (l *logger) Log(v interface{}) {
	log.Print(v)
}

type IAdd interface {
	Add(x int, y int) int
}

type add struct {
	ferry.Feature
}

func (a *add) Add(x int, y int) int {
	return (x + y)
}

type ILogic interface {
}

type logic struct {
	ferry.Feature
	t  *testing.T
	wg *sync.WaitGroup
}

func (l *logic) OnInit(sandbox ferry.ISandbox) {
	l.Feature.OnInit(sandbox)
	l.Book("IAdd")
	l.Book("ILogger")
}

func (l *logic) OnStart() {
	l.Call("ILogger", "Log", "Hello Ferry!")

	result, err := l.CallWithResult("IAdd", "Add", 1, 2)
	if nil != err {
		l.t.Error(err)
	} else {
		l.t.Log(fmt.Sprintf("Add result is: %d\n", result[0]))
	}
	l.wg.Done()
}

func TestOneDock(t *testing.T) {
	network.Mock("tcp")

	var wg sync.WaitGroup
	wg.Add(1)

	go ferry.Serve("127.0.0.1:55555")

	go ferry.Startup("127.0.0.1:55555", "1",
		ferry.Carry("ILogger", &logger{}, true),
		ferry.Carry("IAdd", &add{}, true),
		ferry.Carry("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	ferry.Close()
}

func TestMultiDocks(t *testing.T) {
	network.Mock("tcp")

	var wg sync.WaitGroup
	wg.Add(1)

	go ferry.Serve("127.0.0.1:55555")

	go ferry.Startup("127.0.0.1:55555", "1",
		ferry.Carry("ILogger", &logger{}, true))

	go ferry.Startup("127.0.0.1:55555", "2",
		ferry.Carry("IAdd", &add{}, true))

	go ferry.Startup("127.0.0.1:55555", "3",
		ferry.Carry("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	ferry.Close()
}
