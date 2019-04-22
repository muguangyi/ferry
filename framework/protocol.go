// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

const (
	ERROR             uint = 0
	REGISTER_REQUEST  uint = 1
	REGISTER_RESPONSE uint = 2
	IMPORT_REQUEST    uint = 3
	IMPORT_RESPONSE   uint = 4
	QUERY_REQUEST     uint = 5
	QUERY_RESPONSE    uint = 6
	RPC_REQUEST       uint = 7
	RPC_RESPONSE      uint = 8
)

func protoMaker(id uint) interface{} {
	switch id {
	case ERROR:
		return new(protoError)
	case REGISTER_REQUEST:
		return new(protoRegisterRequest)
	case REGISTER_RESPONSE:
		return new(protoRegisterResponse)
	}

	return nil
}

type protoError struct {
	error string
}

type protoRegisterRequest struct {
	units []string
}

type protoRegisterResponse struct {
	port int
}

type protoImportRequest struct {
	units []string
}

type protoImportResponse struct {
	unions []string
}

type protoQueryRequest struct {
	unit string
}

type protoQueryResponse struct {
	unionAddr string
}

type protoRpcRequest struct {
	unitName string
	method   string
	args     []interface{}
}

type protoRpcResponse struct {
	unitNmae string
	method   string
	result   interface{}
}
