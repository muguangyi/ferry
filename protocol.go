// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"io"

	"github.com/muguangyi/ferry/codec"
)

type cProtoType uint8

const (
	cError            cProtoType = 0x0 // Error
	cHeartbeat        cProtoType = 0x1 // Heartbeat
	cReady            cProtoType = 0x2 // Ready
	cRegisterRequest  cProtoType = 0x3 // Register request
	cRegisterResponse cProtoType = 0x4 // Register response
	cQueryRequest     cProtoType = 0x5 // Query request
	cQueryResponse    cProtoType = 0x6 // Query response
	cRpcRequest       cProtoType = 0x7 // RPC request
	cRpcResponse      cProtoType = 0x8 // RPC response
)

func protoMaker(id cProtoType) IProto {
	switch id {
	case cError:
		return new(protoError)
	case cHeartbeat:
		return new(protoHeartbeat)
	case cRegisterRequest:
		return new(protoRegisterRequest)
	case cRegisterResponse:
		return new(protoRegisterResponse)
	case cReady:
		return new(protoReady)
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

// Error
type protoError struct {
	Error string
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

// Heartbeat
type protoHeartbeat struct {
}

func (p *protoHeartbeat) Marshal(writer io.Writer) error {
	return nil
}

func (p *protoHeartbeat) Unmarshal(reader io.Reader) error {
	return nil
}

// Ready.
type protoReady struct {
	Slots []string
}

func (p *protoReady) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Slots).Encode(writer)
}

func (p *protoReady) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	arr, err := any.Arr()
	if nil != err {
		return err
	}

	p.Slots = make([]string, len(arr))
	for i, iv := range arr {
		p.Slots[i] = iv.(string)
	}

	return nil
}

// Register request
type protoRegisterRequest struct {
	Slots []string
}

func (p *protoRegisterRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Slots).Encode(writer)
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

	p.Slots = make([]string, len(arr))
	for i, iv := range arr {
		p.Slots[i] = iv.(string)
	}

	return nil
}

// Register response
type protoRegisterResponse struct {
	Port int
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

// Query request
type protoQueryRequest struct {
	Slot string
}

func (p *protoQueryRequest) Marshal(writer io.Writer) error {
	return codec.NewAny(p.Slot).Encode(writer)
}

func (p *protoQueryRequest) Unmarshal(reader io.Reader) error {
	any := codec.NewAny(nil)
	err := any.Decode(reader)
	if nil != err {
		return err
	}

	p.Slot, err = any.String()
	return err
}

// Query response
type protoQueryResponse struct {
	DockAddr string
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

// RPC request
type protoRpcRequest struct {
	Index      int64
	Slot       string
	Method     string
	Args       []interface{}
	WithResult bool
}

func (p *protoRpcRequest) Marshal(writer io.Writer) error {
	err := codec.NewAny(p.Index).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Slot).Encode(writer)
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
	p.Slot, err = any.String()
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

// RPC response
type protoRpcResponse struct {
	Index  int64
	Slot   string
	Method string
	Result []interface{}
	Err    string
}

func (p *protoRpcResponse) Marshal(writer io.Writer) error {
	err := codec.NewAny(p.Index).Encode(writer)
	if nil != err {
		return err
	}

	err = codec.NewAny(p.Slot).Encode(writer)
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
	p.Slot, err = any.String()
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
