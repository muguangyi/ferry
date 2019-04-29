// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/muguangyi/seek/chancall"
)

func TestName(t *testing.T) {
	callee := chancall.NewCallee("target", new(targetObject))
	if "target" != callee.Name() {
		t.Fail()
	}
}

type targetObject struct {
}

func (t targetObject) F0() {
	fmt.Println("f0 executed!")
}

func (t targetObject) F1() int {
	return 1
}

func (t targetObject) Add(x int, y int) int {
	return x + y
}

func (t targetObject) LongWaitCall() {
	time.Sleep(2 * time.Second)
}

func TestNormalCall(t *testing.T) {
	callee := chancall.NewCallee("target", new(targetObject))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		caller := chancall.NewCaller(callee)

		err := caller.Call("F0")
		if nil != err {
			fmt.Println(err)
		}

		r1, err := caller.CallWithResult("F1")
		if nil != err {
			fmt.Println(err)
		} else {
			fmt.Println(r1)
		}

		r2, err := caller.CallWithResult("Add", 1, 2)
		if nil != err {
			fmt.Println(err)
		} else {
			fmt.Println(r2)
		}

		wg.Done()
	}()

	wg.Wait()
}

func TestTimeout(t *testing.T) {
	callee := chancall.NewCallee("target", new(targetObject))
	callee.SetTimeout("LongWaitCall", 0.1)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		caller := chancall.NewCaller(callee)

		err := caller.Call("LongWaitCall")
		if nil != err {
			fmt.Println(err)
		} else {
			t.Fail()
		}

		wg.Done()
	}()

	wg.Wait()
}
