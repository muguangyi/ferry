// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ship

const (
	cError            uint = 0
	cHeartbeat        uint = 1
	cRegisterRequest  uint = 2
	cRegisterResponse uint = 3
	cImportRequest    uint = 4
	cImportResponse   uint = 5
	cQueryRequest     uint = 6
	cQueryResponse    uint = 7
	cRpcRequest       uint = 8
	cRpcResponse      uint = 9
)

func protoMaker(id uint) interface{} {
	switch id {
	case cError:
		return new(protoError)
	case cHeartbeat:
		return new(protoHeartbeat)
	case cRegisterRequest:
		return new(protoRegisterRequest)
	case cRegisterResponse:
		return new(protoRegisterResponse)
	case cImportRequest:
		return new(protoImportRequest)
	case cImportResponse:
		return new(protoImportResponse)
	case cQueryRequest:
		return new(protoQueryRequest)
	case cQueryResponse:
		return new(protoQueryResponse)
	case cRpcRequest:
		return new(protoRpcRequest)
	case cRpcResponse:
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
	Signalers []string `json:"signalers"`
}

type protoRegisterResponse struct {
	Port int `json:"port"`
}

type protoImportRequest struct {
	Signalers []string `json:"signalers"`
}

type protoImportResponse struct {
	Unions []string `json:"unions"`
}

type protoQueryRequest struct {
	Signaler string `json:"signaler"`
}

type protoQueryResponse struct {
	UnionAddr string `json:"union-addr"`
}

type protoRpcRequest struct {
	Index      int64         `json:"index"`
	SignalerId string        `json:"signaler-id"`
	Method     string        `json:"method"`
	Args       []interface{} `json:"args"`
	WithResult bool          `json:"with-result"`
}

type protoRpcResponse struct {
	Index      int64         `json:"index"`
	SignalerId string        `json:"signaler-id"`
	Method     string        `json:"method"`
	Result     []interface{} `json:"result"`
	Err        string        `json:"error"`
}
