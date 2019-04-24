// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/muguangyi/unite/chancall"
)

func Test(t *testing.T) {
	callee := chancall.NewCallee()

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
	}()

	wg.Wait()
	wg.Add(1)

	go func() {
		caller := chancall.NewCaller(callee)

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
