// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/muguangyi/gounite/rpc"
)

func Test(t *testing.T) {
	callee := rpc.NewCallee()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		callee.Bind("f0", func(args []interface{}) {
			fmt.Println("f0 executed!")
		})

		callee.Bind("f1", func(args []interface{}) interface{} {
			return 1
		})

		callee.Bind("add", func(args []interface{}) interface{} {
			return args[0].(int) + args[1].(int)
		})

		wg.Done()

		callee.Handling()
	}()

	wg.Wait()
	wg.Add(1)

	go func() {
		caller := rpc.NewCaller(callee)

		err := caller.Call("f0")
		if nil != err {
			fmt.Println(err)
		}

		r1, err := caller.CallWithResult("f1")
		if nil != err {
			fmt.Println(err)
		} else {
			fmt.Println(r1)
		}

		r2, err := caller.CallWithResult("add", 1, 2)
		if nil != err {
			fmt.Println(err)
		} else {
			fmt.Println(r2)
		}

		wg.Done()
	}()

	wg.Wait()
}
