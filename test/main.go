// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/muguangyi/ferry"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go ferry.Serve("127.0.0.1:9999")

	go ferry.Startup("127.0.0.1:9999", "util",
		ferry.Carry("IMath", newMath(), true))

	go ferry.Startup("127.0.0.1:9999", "logic",
		ferry.Carry("IGame", newGame(&wg), true),
		ferry.Carry("ILobby", newLobby(&wg), true))

	wg.Wait()
	ferry.Close()
	fmt.Println("Completed!")
}

//go:generate ferry.gen
