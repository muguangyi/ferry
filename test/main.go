// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/muguangyi/unite/unite"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	unite.RunHub("127.0.0.1:9999")

	unite.Run("127.0.0.1:9999", "util",
		unite.NewUnit("IMath", newMath(), true))

	unite.Run("127.0.0.1:9999", "logic",
		unite.NewUnit("IGame", newGame(&wg), true),
		unite.NewUnit("ILobby", newLobby(&wg), true))

	wg.Wait()
	fmt.Println("Completed!")
}

//go:generate unite.gen
