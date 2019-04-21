// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package internal

type proto struct {
	id uint
}

type protoError struct {
	proto
	error string
}

type protoRegisterUnionRequest struct {
	proto
	units []string
}

type protoRegisterUnionResponse struct {
	proto
	port int
}

type protoQueryUnitRequest struct {
	proto
	unit string
}

type protoQueryUnitResponse struct {
	proto
	unionAddr string
}

type protoBindUnionRequest struct {
	proto
	units []string
}

type protoBindUnionResponse struct {
	proto
}

type protoRPC struct {
	proto
	unitName string
	method   string
	args     []interface{}
}
