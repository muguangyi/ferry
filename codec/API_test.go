// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec_test

import (
	"bytes"
	"testing"

	"github.com/muguangyi/ferry/codec"
)

func Test_String(t *testing.T) {
	buf := &bytes.Buffer{}
	err := codec.NewAny("ferry").Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.String()
	if err != nil {
		t.Error(err)
	}

	if v != "ferry" {
		t.Fail()
	}
}

func Test_Bool(t *testing.T) {
	buf := &bytes.Buffer{}

	err := codec.NewAny(true).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Bool()
	if err != nil {
		t.Error(err)
	}
	if v != true {
		t.Fail()
	}

	err = codec.NewAny(false).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err = any.Bool()
	if err != nil {
		t.Error(err)
	}
	if v != false {
		t.Fail()
	}
}

func Test_Nil(t *testing.T) {
	var value interface{} = nil

	buf := &bytes.Buffer{}

	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(buf.Bytes(), []byte{0xc0}) != 0 {
		t.Fail()
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if any.Any() != nil {
		t.Fail()
	}
}

func Test_Uint8(t *testing.T) {
	var value uint8 = 255

	buf := &bytes.Buffer{}

	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Uint8()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Uint16(t *testing.T) {
	var value uint16 = 520

	buf := &bytes.Buffer{}

	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Uint16()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Uint32(t *testing.T) {
	var value uint32 = 520

	buf := &bytes.Buffer{}
	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Uint32()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Uint64(t *testing.T) {
	var value uint64 = 520

	buf := &bytes.Buffer{}
	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Uint64()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Int8(t *testing.T) {
	var value int8 = -127

	buf := &bytes.Buffer{}
	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Int8()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Int16(t *testing.T) {
	var value int16 = -520

	buf := &bytes.Buffer{}
	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Int16()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Int32(t *testing.T) {
	var value int32 = -520

	buf := &bytes.Buffer{}
	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Int32()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Int64(t *testing.T) {
	var value int64 = -520

	buf := &bytes.Buffer{}
	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Int64()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Float32(t *testing.T) {
	var value float32 = 0.520

	buf := &bytes.Buffer{}
	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Float32()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Float64(t *testing.T) {
	var value float64 = -0.00520

	buf := &bytes.Buffer{}

	err := codec.NewAny(value).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	v, err := any.Float64()
	if err != nil {
		t.Error(err)
	}

	if v != value {
		t.Fail()
	}
}

func Test_Array(t *testing.T) {
	values := []int{0, 1, 2, 3, 4, 5, 6}

	buf := &bytes.Buffer{}

	err := codec.NewAny(values).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	arr, err := any.Arr()
	if err != nil {
		t.Error(err)
	}

	for i, n := range values {
		if n != int(arr[i].(int8)) {
			t.Fail()
		}
	}
}

func Test_Map(t *testing.T) {
	values := map[string]int{
		"a": 1,
		"b": 2,
	}

	buf := &bytes.Buffer{}

	err := codec.NewAny(values).Encode(buf)
	if err != nil {
		t.Error(err)
	}

	any := codec.NewAny(nil)
	err = any.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	dict, err := any.Map()
	if err != nil {
		t.Error(err)
	}

	for k, n := range dict {
		s, ok := values[k.(string)]
		if !ok || s != int(n.(int8)) {
			t.Fail()
		}
	}
}
