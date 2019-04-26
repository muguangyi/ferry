// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unite

var registry map[string]interface{} = make(map[string]interface{})

func register(id string, maker interface{}) {
	registry[id] = maker
}

func tryMake(id string, u IUnit) (interface{}, bool) {
	maker := registry[id]
	if nil != maker {
		return maker.(func(u IUnit) interface{})(u), true
	}

	return nil, false
}
