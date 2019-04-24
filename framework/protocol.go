// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package framework

// import (
// 	"encoding/json"
// )

const (
	ERROR             uint = 0
	HEARTBEAT         uint = 1
	REGISTER_REQUEST  uint = 2
	REGISTER_RESPONSE uint = 3
	IMPORT_REQUEST    uint = 4
	IMPORT_RESPONSE   uint = 5
	QUERY_REQUEST     uint = 6
	QUERY_RESPONSE    uint = 7
	RPC_REQUEST       uint = 8
	RPC_RESPONSE      uint = 9
)

func protoMaker(id uint) interface{} {
	switch id {
	case ERROR:
		return new(protoError)
	case HEARTBEAT:
		return new(protoHeartbeat)
	case REGISTER_REQUEST:
		return new(protoRegisterRequest)
	case REGISTER_RESPONSE:
		return new(protoRegisterResponse)
	case IMPORT_REQUEST:
		return new(protoImportRequest)
	case IMPORT_RESPONSE:
		return new(protoImportResponse)
	case QUERY_REQUEST:
		return new(protoQueryRequest)
	case QUERY_RESPONSE:
		return new(protoQueryResponse)
	case RPC_REQUEST:
		return new(protoRpcRequest)
	case RPC_RESPONSE:
		return new(protoRpcResponse)
	}

	return nil
}

type protoError struct {
	Error string `json:"error"`
}

type protoHeartbeat struct {
}

type protoRegisterRequest struct {
	Units []string `json:"units"`
}

type protoRegisterResponse struct {
	Port int `json:"port"`
}

type protoImportRequest struct {
	Units []string `json:"units"`
}

type protoImportResponse struct {
	Unions []string `json:"unions"`
}

type protoQueryRequest struct {
	Unit string `json:"unit"`
}

type protoQueryResponse struct {
	UnionAddr string `json:"union-addr"`
}

type protoRpcRequest struct {
	Index      int64         `json:"index"`
	UnitId     string        `json:"unit-id"`
	Method     string        `json:"method"`
	Args       []interface{} `json:"args"`
	WithResult bool          `json:"with-result"`
}

type protoRpcResponse struct {
	Index  int64       `json:"index"`
	UnitId string      `json:"unit-id"`
	Method string      `json:"method"`
	Result interface{} `json:"result"`
}
