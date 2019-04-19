// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

import (
	"bytes"
	"encoding/json"
)

type jsonSerializer struct {
}

func (j *jsonSerializer) Marshal(obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if nil != err {
		panic("JSON marshal failed!")
	}
	length := len(data)

	header := make([]byte, 4)
	header[0] = byte(length)
	header[0] = byte(length >> 8)
	header[0] = byte(length >> 16)
	header[0] = byte(length >> 24)

	return joinBytes(header, data)
}

func (j *jsonSerializer) Unmarshal(data []byte) interface{} {
	length := (int(data[0]) | int(data[1])<<8 | int(data[2])<<16 | int(data[3])<<24)

	var obj interface{}
	err := json.Unmarshal(data[4:(4+length)], obj)
	if nil != err {
		panic("JSON unmarshal failed!")
	}

	return obj
}

func (j *jsonSerializer) Slice(source []byte) int {
	if len(source) < 4 {
		return 0
	}

	length := (int(source[0]) | int(source[1])<<8 | int(source[2])<<16 | int(source[3])<<24)
	if len(source) < (4 + length) {
		return 0
	}

	return (4 + length)
}

func joinBytes(data ...[]byte) []byte {
	return bytes.Join(data, []byte(""))
}
