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
	_, err := codec.Encode(buf, "ferry")
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v.(string) != "ferry" {
		t.Fail()
	}
}

func Test_Bool(t *testing.T) {
	buf := &bytes.Buffer{}

	_, err := codec.Encode(buf, true)
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v.(bool) != true {
		t.Fail()
	}

	_, err = codec.Encode(buf, false)
	if err != nil {
		t.Error(err)
	}

	v, err = codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v.(bool) != false {
		t.Fail()
	}
}

func Test_Nil(t *testing.T) {
	var value interface{} = nil

	buf := &bytes.Buffer{}

	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(buf.Bytes(), []byte{0xc0}) != 0 {
		t.Fail()
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v != nil {
		t.Fail()
	}
}

func Test_Uint8(t *testing.T) {
	var value uint8 = 255

	buf := &bytes.Buffer{}

	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v.(uint8) != value {
		t.Fail()
	}
}

func Test_Uint16(t *testing.T) {
	var value uint16 = 520

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v.(uint16) != value {
		t.Fail()
	}
}

func Test_Uint32(t *testing.T) {
	var value uint32 = 520

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Fail()
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Fail()
	}

	if uint32(v.(uint16)) != value {
		t.Fail()
	}
}

func Test_Uint64(t *testing.T) {
	var value uint64 = 520

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Fail()
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Fail()
	}

	if uint64(v.(uint16)) != value {
		t.Fail()
	}
}

func Test_Int8(t *testing.T) {
	var value int8 = -127

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Fail()
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Fail()
	}

	if v.(int8) != value {
		t.Fail()
	}
}

func Test_Int16(t *testing.T) {
	var value int16 = -520

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Fail()
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Fail()
	}

	if v.(int16) != value {
		t.Fail()
	}
}

func Test_Int32(t *testing.T) {
	var value int32 = -520

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Fail()
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Fail()
	}

	if int32(v.(int16)) != value {
		t.Fail()
	}
}

func Test_Int64(t *testing.T) {
	var value int64 = -520

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Fail()
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Fail()
	}

	if int64(v.(int16)) != value {
		t.Fail()
	}
}

func Test_Float32(t *testing.T) {
	var value float32 = 0.520

	buf := &bytes.Buffer{}
	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v.(float32) != value {
		t.Fail()
	}
}

func Test_Float64(t *testing.T) {
	var value float64 = -0.00520

	buf := &bytes.Buffer{}

	_, err := codec.Encode(buf, value)
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	if v.(float64) != value {
		t.Fail()
	}
}

func Test_Array(t *testing.T) {
	values := []int{0, 1, 2, 3, 4, 5, 6}

	buf := &bytes.Buffer{}

	_, err := codec.Encode(buf, values)
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	arr, ok := v.([]interface{})
	if !ok {
		t.Fail()
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

	_, err := codec.Encode(buf, values)
	if err != nil {
		t.Error(err)
	}

	v, err := codec.Decode(buf)
	if err != nil {
		t.Error(err)
	}

	dict, ok := v.(map[interface{}]interface{})
	if !ok {
		t.Fail()
	}

	for k, n := range dict {
		s, ok := values[k.(string)]
		if !ok || s != int(n.(int8)) {
			t.Fail()
		}
	}
}
