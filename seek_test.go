// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek_test

import (
	"log"
	"sync"
	"testing"

	"github.com/muguangyi/seek"
	"github.com/muguangyi/seek/network"
)

type ILogger interface {
	Log(v interface{})
}

type logger struct {
	seek.Signal
}

func (l *logger) Log(v interface{}) {
	log.Print(v)
}

type IAdd interface {
	Add(x int, y int) int
}

type add struct {
	seek.Signal
}

func (a *add) Add(x int, y int) int {
	return (x + y)
}

type ILogic interface {
}

type logic struct {
	seek.Signal
	t  *testing.T
	wg *sync.WaitGroup
}

func (l *logic) OnInit(signaler seek.ISignaler) {
	l.Signal.OnInit(signaler)
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

	go seek.Serve("127.0.0.1:55555")

	go seek.Startup("127.0.0.1:55555", "1",
		seek.NewSignaler("ILogger", &logger{}, true),
		seek.NewSignaler("IAdd", &add{}, true),
		seek.NewSignaler("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	seek.Close()
}

func TestMultiUnions(t *testing.T) {
	network.Mock(true)

	var wg sync.WaitGroup
	wg.Add(1)

	go seek.Serve("127.0.0.1:55555")

	go seek.Startup("127.0.0.1:55555", "1",
		seek.NewSignaler("ILogger", &logger{}, true))

	go seek.Startup("127.0.0.1:55555", "2",
		seek.NewSignaler("IAdd", &add{}, true))

	go seek.Startup("127.0.0.1:55555", "3",
		seek.NewSignaler("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	seek.Close()
}
