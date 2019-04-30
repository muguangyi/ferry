// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek_test

import (
	"sync"
	"testing"

	"github.com/muguangyi/seek/network"
	"github.com/muguangyi/seek/seek"
)

type IAdd interface {
	Do(x int, y int) int
}

type add struct {
	seek.Signal
}

func (a *add) Do(x int, y int) int {
	return (x + y)
}

func (a *add) OnInit(signaler seek.ISignaler) {
	a.Signal.OnInit(signaler)
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
	signaler.Book("IAdd")
}

func (l *logic) OnStart() {
	result, err := l.CallWithResult("IAdd", "Do", 1, 2)
	if nil != err {
		l.t.Fail()
	} else if result[0].(int) != 3 {
		l.t.Fail()
	}
	l.wg.Done()
}

func TestOneUnion(t *testing.T) {
	network.Mock(true)

	var wg sync.WaitGroup
	wg.Add(1)

	seek.Serve("127.0.0.1:55555")

	seek.Startup("127.0.0.1:55555", "1",
		seek.NewSignaler("IAdd", &add{}, true),
		seek.NewSignaler("ILogic", &logic{t: t, wg: &wg}, true))

	wg.Wait()
}
