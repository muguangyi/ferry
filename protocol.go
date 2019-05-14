// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

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
	Sandboxes []string `json:"sandboxes"`
}

type protoRegisterResponse struct {
	Port int `json:"port"`
}

type protoImportRequest struct {
	Sandboxes []string `json:"sandboxes"`
}

type protoImportResponse struct {
	Docks []string `json:"docks"`
}

type protoQueryRequest struct {
	Sandbox string `json:"sandbox"`
}

type protoQueryResponse struct {
	DockAddr string `json:"dock-addr"`
}

type protoRpcRequest struct {
	Index      int64         `json:"index"`
	SandboxId  string        `json:"sandbox-id"`
	Method     string        `json:"method"`
	Args       []interface{} `json:"args"`
	WithResult bool          `json:"with-result"`
}

type protoRpcResponse struct {
	Index     int64         `json:"index"`
	SandboxId string        `json:"sandbox-id"`
	Method    string        `json:"method"`
	Result    []interface{} `json:"result"`
	Err       string        `json:"error"`
}
