// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"io"

	"github.com/muguangyi/ferry/codec"
)

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

func protoMaker(id uint) IProto {
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

func (p *protoError) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Error).Encode(writer)
}

func (p *protoError) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Error, err = any.String()
	return err
}

type protoHeartbeat struct {
}

func (p *protoHeartbeat) Marshal(writer io.Writer) error {
	return nil
}

func (p *protoHeartbeat) Unmarshal(reader io.Reader) error {
	return nil
}

type protoRegisterRequest struct {
	Sandboxes []string `json:"sandboxes"`
}

func (p *protoRegisterRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Sandboxes).Encode(writer)
}

func (p *protoRegisterRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	arr, err := any.Arr()
	if nil != err {
		return err
	}

	p.Sandboxes = make([]string, len(arr))
	for i, iv := range arr {
		p.Sandboxes[i] = iv.(string)
	}

	return nil
}

type protoRegisterResponse struct {
	Port int `json:"port"`
}

func (p *protoRegisterResponse) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Port).Encode(writer)
}

func (p *protoRegisterResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Port, err = any.Int()
	return err
}

type protoImportRequest struct {
	Sandboxes []string `json:"sandboxes"`
}

func (p *protoImportRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Sandboxes).Encode(writer)
}

func (p *protoImportRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	arr, err := any.Arr()
	if nil != err {
		return err
	}

	p.Sandboxes = make([]string, len(arr))
	for i, iv := range arr {
		p.Sandboxes[i] = iv.(string)
	}

	return nil
}

type protoImportResponse struct {
	Docks []string `json:"docks"`
}

func (p *protoImportResponse) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Docks).Encode(writer)
}

func (p *protoImportResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	arr, err := any.Arr()
	if nil != err {
		return err
	}

	p.Docks = make([]string, len(arr))
	for i, iv := range arr {
		p.Docks[i] = iv.(string)
	}

	return nil
}

type protoQueryRequest struct {
	Sandbox string `json:"sandbox"`
}

func (p *protoQueryRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Sandbox).Encode(writer)
}

func (p *protoQueryRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Sandbox, err = any.String()
	return err
}

type protoQueryResponse struct {
	DockAddr string `json:"dock-addr"`
}

func (p *protoQueryResponse) Marshal(writer io.Writer) error {
	return codec.NewAny(p.DockAddr).Encode(writer)
}

func (p *protoQueryResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.DockAddr, err = any.String()
	return err
}

type protoRpcRequest struct {
	Index      int64         `json:"index"`
	SandboxId  string        `json:"sandbox-id"`
	Method     string        `json:"method"`
	Args       []interface{} `json:"args"`
	WithResult bool          `json:"with-result"`
}

func (p *protoRpcRequest) Marshal(writer io.Writer) error {
	err := codec.NewAny(p.Index).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.SandboxId).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Method).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Args).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.WithResult).Encode(writer)
	if nil != err {
		return err
	}

	return nil
}

func (p *protoRpcRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)

	err := any.Decode(reader)
	if nil != err {
		return err
	}
	p.Index, err = any.Int64()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.SandboxId, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Method, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Args, err = any.Arr()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.WithResult, err = any.Bool()
	if nil != err {
		return err
	}

	return nil
}

type protoRpcResponse struct {
	Index     int64         `json:"index"`
	SandboxId string        `json:"sandbox-id"`
	Method    string        `json:"method"`
	Result    []interface{} `json:"result"`
	Err       string        `json:"error"`
}

func (p *protoRpcResponse) Marshal(writer io.Writer) error {
	err := codec.NewAny(p.Index).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.SandboxId).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Method).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Result).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Err).Encode(writer)
	if nil != err {
		return err
	}

	return nil
}

func (p *protoRpcResponse) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)

	err := any.Decode(reader)
	if nil != err {
		return err
	}
	p.Index, err = any.Int64()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.SandboxId, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Method, err = any.String()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Result, err = any.Arr()
	if nil != err {
		return err
	}

	err = any.Decode(reader)
	if nil != err {
		return err
	}
	p.Err, err = any.String()
	if nil != err {
		return err
	}

	return nil
}
