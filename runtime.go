// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"log"
	"os"
	"os/signal"
)

type instance interface {
	Close()
}

var insts []instance = make([]instance, 0)

func wait(inst instance) {
	insts = append(insts, inst)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c

	log.Printf("instance closed (signal: %v)", sig)
}

func destroy() {
	for i := len(insts) - 1; i >= 0; i-- {
		insts[i].Close()
	}
	insts = insts[:0]
}
