// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec

// Subset format from MsgPack (https://github.com/msgpack/msgpack/blob/master/spec.md).
const (
	cNil byte = 0xc0

	cFalse byte = 0xc2
	cTrue  byte = 0xc3

	cBin8  byte = 0xc4
	cBin16 byte = 0xc5
	cBin32 byte = 0xc6

	cFloat32 byte = 0xca
	cFloat64 byte = 0xcb

	cUint8  byte = 0xcc
	cUint16 byte = 0xcd
	cUint32 byte = 0xce
	cUint64 byte = 0xcf

	cInt8  byte = 0xd0
	cInt16 byte = 0xd1
	cInt32 byte = 0xd2
	cInt64 byte = 0xd3

	cStr8  byte = 0xd9
	cStr16 byte = 0xda
	cStr32 byte = 0xdb

	cArr16 byte = 0xdc
	cArr32 byte = 0xdd

	cMap16 byte = 0xde
	cMap32 byte = 0xdf
)

const (
	cInt32Size = 4
	cInt64Size = 8

	cUint8Max  = 2 << (8 - 1)
	cUint16Max = 2 << (16 - 1)
	cUint32Max = 2 << (32 - 1)

	cInt8Max  = 2 << (8 - 2)
	cInt8Min  = -cInt8Max
	cInt16Max = 2 << (16 - 2)
	cInt16Min = -cInt16Max
	cInt32Max = 2 << (32 - 2)
	cInt32Min = -cInt32Max
)

type bytes []byte
type bytes1 [1]byte
type bytes2 [2]byte
type bytes4 [4]byte
type bytes8 [8]byte
