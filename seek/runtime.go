// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

type instance interface {
	Close()
}

var insts []instance = make([]instance, 0)

func watch(inst instance) {
	insts = append(insts, inst)
}

func destroy() {
	for i := len(insts) - 1; i >= 0; i-- {
		insts[i].Close()
	}
	insts = insts[:0]
}
