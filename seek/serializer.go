// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package seek

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/muguangyi/seek/network"
)

type packer struct {
	Id uint        `json:"id"`
	P  interface{} `json:"p"`
}

type unpacker struct {
	Id uint            `json:"id"`
	P  json.RawMessage `json:"p"`
}

func newSerializer() network.ISerializer {
	return &serializer{
		maker: protoMaker,
	}
}

type serializer struct {
	maker func(id uint) interface{}
}

func (s *serializer) Marshal(obj interface{}) []byte {
	switch obj.(type) {
	case *packer:
		data, err := json.Marshal(obj)
		if nil != err {
			log.Fatal(err)
			return nil
		}
		length := len(data)

		header := make([]byte, 4)
		header[0] = byte(length)
		header[1] = byte(length >> 8)
		header[2] = byte(length >> 16)
		header[3] = byte(length >> 24)

		return joinBytes(header, data)
	default:
		log.Fatal("Unknown type!")
	}

	return nil
}

func (s *serializer) Unmarshal(data []byte) interface{} {
	length := (int(data[0]) | int(data[1])<<8 | int(data[2])<<16 | int(data[3])<<24)
	body := data[4 : 4+length]
	obj := &unpacker{}
	err := json.Unmarshal(body, obj)
	if nil != err {
		log.Fatal(err)
		return nil
	}

	p := s.maker(obj.Id)
	err = json.Unmarshal(obj.P, p)
	if nil != err {
		log.Fatal(err)
		return nil
	}

	return &packer{
		Id: obj.Id,
		P:  p,
	}
}

func (s *serializer) Slice(source []byte) int {
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
