// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ship_test

import (
	"log"
	"sync"
	"testing"

	"github.com/muguangyi/ship"
	"github.com/muguangyi/ship/network"
)

type ILogger interface {
	Log(v interface{})
}

type logger struct {
	ship.Feature
}

func (l *logger) Log(v interface{}) {
	log.Print(v)
}

type IAdd interface {
	Add(x int, y int) int
}

type add struct {
	seek.Feature
}

func (a *add) Add(x int, y int) int {
	return (x + y)
}

type ILogic interface {
}

type logic struct {
	ship.Feature
	t  *testing.T
	wg *sync.WaitGroup
}

func (l *logic) OnInit(sandbox ship.ISandbox) {
	l.Feature.OnInit(sandbox)
	l.Book("IAdd")
	l.Book("ILogger")
}

func (l *logic) OnStart() {
	l.Call("ILogger", "Log", "Hello Seek!")

	result, err := l.CallWithResult("IAdd", "Add", 1, 2)
	if nil != err {
		l.t.Fail()
	} else {
		l.t.Log(result[0])
	}
	l.wg.Done()
}

func TestOneUnion(t *testing.T) {
	network.Mock(true)

	var wg sync.WaitGroup
	wg.Add(1)

	go ship.Serve("127.0.0.1:55555")

	go ship.Startup("127.0.0.1:55555", "1",
		ship.Carry("ILogger", &logger{}, true),
		ship.Carry("IAdd", &add{}, true),
		ship.Carry("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	ship.Close()
}

func TestMultiUnions(t *testing.T) {
	network.Mock(true)

	var wg sync.WaitGroup
	wg.Add(1)

	go ship.Serve("127.0.0.1:55555")

	go ship.Startup("127.0.0.1:55555", "1",
		ship.Carry("ILogger", &logger{}, true))

	go ship.Startup("127.0.0.1:55555", "2",
		ship.Carry("IAdd", &add{}, true))

	go ship.Startup("127.0.0.1:55555", "3",
		ship.Carry("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	ship.Close()
}
