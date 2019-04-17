// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package rpc

type callRequest struct {
	function     interface{}
	args         []interface{}
	callResponse chan *callResponse
}

type callResponse struct {
	result interface{}
	err    error
}
