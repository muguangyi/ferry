// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

const (
	cERROR             uint = 0
	cHEARTBEAT         uint = 1
	cREGISTER_REQUEST  uint = 2
	cREGISTER_RESPONSE uint = 3
	cIMPORT_REQUEST    uint = 4
	cIMPORT_RESPONSE   uint = 5
	cQUERY_REQUEST     uint = 6
	cQUERY_RESPONSE    uint = 7
	cRPC_REQUEST       uint = 8
	cRPC_RESPONSE      uint = 9
)

func protoMaker(id uint) interface{} {
	switch id {
	case cERROR:
		return new(protoError)
	case cHEARTBEAT:
		return new(protoHeartbeat)
	case cREGISTER_REQUEST:
		return new(protoRegisterRequest)
	case cREGISTER_RESPONSE:
		return new(protoRegisterResponse)
	case cIMPORT_REQUEST:
		return new(protoImportRequest)
	case cIMPORT_RESPONSE:
		return new(protoImportResponse)
	case cQUERY_REQUEST:
		return new(protoQueryRequest)
	case cQUERY_RESPONSE:
		return new(protoQueryResponse)
	case cRPC_REQUEST:
		return new(protoRpcRequest)
	case cRPC_RESPONSE:
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
