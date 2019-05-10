// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/muguangyi/ship"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go ship.Serve("127.0.0.1:9999")

	go ship.Startup("127.0.0.1:9999", "util",
		ship.Carry("IMath", newMath(), true))

	go ship.Startup("127.0.0.1:9999", "logic",
		ship.Carry("IGame", newGame(&wg), true),
		ship.Carry("ILobby", newLobby(&wg), true))

	wg.Wait()
	ship.Close()
	fmt.Println("Completed!")
}

//go:generate ship.gen
