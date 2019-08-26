// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

var registry map[string]interface{} = make(map[string]interface{})

func register(name string, maker interface{}) {
	registry[name] = maker
}

func tryMake(name string, s ISlot) (interface{}, bool) {
	maker := registry[name]
	if nil != maker {
		return maker.(func(slot ISlot) interface{})(s), true
	}

	return nil, false
}
