// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

import (
	"sync"
)

type callRequest struct {
	sync.Mutex
	method       string
	args         []interface{}
	callResponse chan *callResponse
	done         bool
}

type callResponse struct {
	result interface{}
	err    error
}
