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
	wg *sync.WaitGroup
}

func (l *logger) Log(v interface{}) {
	log.Print(v)
	l.wg.Done()
}

type IAdd interface {
	Add(x int, y int) int
}

type add struct {
	ferry.Feature
	wg *sync.WaitGroup
}

func (a *add) Add(x int, y int) int {
	a.wg.Done()
	return (x + y)
}

type ILogic interface {
}

type logic struct {
	ferry.Feature
	t  *testing.T
	wg *sync.WaitGroup
}

func (l *logic) OnInit() []string {
	return []string{"IAdd", "ILogger"}
}

func (l *logic) OnStart(s ferry.ISlot) {
	s.Call("ILogger", "Log", "Hello Ferry!")

	result, err := s.CallWithResult("IAdd", "Add", 1, 2)
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
	wg.Add(3)

	go ferry.Serve("127.0.0.1:55555")

	go ferry.Startup("127.0.0.1:55555", "1",
		ferry.Carry("ILogger", &logger{wg: &wg}, true),
		ferry.Carry("IAdd", &add{wg: &wg}, true),
		ferry.Carry("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	ferry.Close()
}

func TestMultiDocks(t *testing.T) {
	network.Mock("tcp")

	var wg sync.WaitGroup
	wg.Add(3)

	go ferry.Serve("127.0.0.1:55555")

	go ferry.Startup("127.0.0.1:55555", "logger",
		ferry.Carry("ILogger", &logger{wg: &wg}, true))

	go ferry.Startup("127.0.0.1:55555", "add",
		ferry.Carry("IAdd", &add{wg: &wg}, true))

	go ferry.Startup("127.0.0.1:55555", "logic",
		ferry.Carry("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()

	ferry.Close()
}
