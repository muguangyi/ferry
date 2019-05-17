// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"unsafe"
)

func encodeNil(writer io.Writer) (n int, err error) {
	return writer.Write(bytes{cNil})
}

func encodeBool(writer io.Writer, value bool) (n int, err error) {
	if value {
		return writer.Write(bytes{cTrue})
	} else {
		return writer.Write(bytes{cFalse})
	}
}

func encodeFloat32(writer io.Writer, value float32) (n int, err error) {
	return encodeUint32(writer, *(*uint32)(unsafe.Pointer(&value)))
}

func encodeFloat64(writer io.Writer, value float64) (n int, err error) {
	return encodeUint64(writer, *(*uint64)(unsafe.Pointer(&value)))
}

func encodeUint8(writer io.Writer, value uint8) (n int, err error) {
	return writer.Write(bytes{cUint8, byte(value)})
}

func encodeUint16(writer io.Writer, value uint16) (n int, err error) {
	if value < cUint8Max {
		return encodeUint8(writer, uint8(value))
	}

	return writer.Write(bytes{cUint16, byte(value >> 8), byte(value)})
}

func encodeUint32(writer io.Writer, value uint32) (n int, err error) {
	if value < cUint16Max {
		return encodeUint16(writer, uint16(value))
	}

	return writer.Write(bytes{cUint32, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

func encodeUint64(writer io.Writer, value uint64) (n int, err error) {
	if value < cUint32Max {
		return encodeUint32(writer, uint32(value))
	}

	return writer.Write(bytes{cUint64, byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

func encodeUint(writer io.Writer, value uint) (n int, err error) {
	switch unsafe.Sizeof(value) {
	case cInt32Size:
		return encodeUint32(writer, *(*uint32)(unsafe.Pointer(&value)))
	case cInt64Size:
		return encodeUint64(writer, *(*uint64)(unsafe.Pointer(&value)))
	}

	return 0, os.ErrNotExist
}

func encodeInt8(writer io.Writer, value int8) (n int, err error) {
	return writer.Write(bytes{cUint8, byte(value)})
}

func encodeInt16(writer io.Writer, value int16) (n int, err error) {
	if value > cInt8Min && value < cInt8Max {
		return encodeInt8(writer, int8(value))
	}

	return writer.Write(bytes{cUint16, byte(uint16(value) >> 8), byte(value)})
}

func encodeInt32(writer io.Writer, value int32) (n int, err error) {
	if value > cInt16Min && value < cInt16Max {
		return encodeInt16(writer, int16(value))
	}

	return writer.Write(bytes{cUint32, byte(uint32(value) >> 24), byte(uint32(value) >> 16), byte(uint32(value) >> 8), byte(value)})
}

func encodeInt64(writer io.Writer, value int64) (n int, err error) {
	if value > cInt32Min && value < cInt32Max {
		return encodeInt32(writer, int32(value))
	}

	return writer.Write(bytes{cUint64, byte(uint64(value) >> 56), byte(uint64(value) >> 48), byte(uint64(value) >> 40), byte(uint64(value) >> 32), byte(uint64(value) >> 24), byte(uint64(value) >> 16), byte(uint64(value) >> 8), byte(value)})
}

func encodeInt(writer io.Writer, value int) (n int, err error) {
	switch unsafe.Sizeof(value) {
	case cInt32Size:
		return encodeInt32(writer, *(*int32)(unsafe.Pointer(&value)))
	case cInt64Size:
		return encodeInt64(writer, *(*int64)(unsafe.Pointer(&value)))
	}

	return 0, os.ErrNotExist
}

func encodeString(writer io.Writer, value string) (n int, err error) {
	data := bytes(value)
	length := len(data)
	n1, err := writer.Write(bytes{cStr32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
	if nil != err {
		return n1, err
	}

	n2, err := writer.Write(data)
	return n1 + n2, err
}

func encodeArray(writer io.Writer, value reflect.Value) (n int, err error) {
	length := value.Len()
	n, err = writer.Write(bytes{cArr32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
	if nil != err {
		return n, err
	}

	for i := 0; i < length; i++ {
		ni, err := encodeValue(writer, value.Index(i))
		if nil != err {
			return n + ni, err
		}

		n += ni
	}

	return n, nil
}

func encodeMap(writer io.Writer, value reflect.Value) (n int, err error) {
	keys := value.MapKeys()
	length := len(keys)
	n, err = writer.Write(bytes{cMap32, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
	if nil != err {
		return n, err
	}

	for _, k := range keys {
		nk, err := encodeValue(writer, k)
		if nil != err {
			return n + nk, err
		}
		n += nk

		nk, err = encodeValue(writer, value.MapIndex(k))
		if nil != err {
			return n + nk, err
		}
		n += nk
	}

	return n, nil
}

func encodeValue(writer io.Writer, value reflect.Value) (n int, err error) {
	if !value.IsValid() || nil == value.Type() {
		return encodeNil(writer)
	}

	switch v := value; v.Kind() {
	case reflect.Bool:
		return encodeBool(writer, v.Bool())
	case reflect.Float32, reflect.Float64:
		return encodeFloat64(writer, v.Float())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return encodeUint64(writer, v.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt64(writer, v.Int())
	case reflect.String:
		return encodeString(writer, v.String())
	case reflect.Array, reflect.Slice:
		return encodeArray(writer, v)
	case reflect.Map:
		return encodeMap(writer, v)
	case reflect.Interface:
		{
			vv := reflect.ValueOf(v.Interface())
			if reflect.Interface != vv.Kind() {
				return encodeValue(writer, vv)
			}
		}
	}

	return 0, fmt.Errorf("Unsupported type: %s", value.Type().String())
}

func encode(writer io.Writer, value interface{}) (n int, err error) {
	if nil == value {
		return encodeNil(writer)
	}

	switch v := value.(type) {
	case bool:
		return encodeBool(writer, v)
	case float32:
		return encodeFloat32(writer, v)
	case float64:
		return encodeFloat64(writer, v)
	case uint8:
		return encodeUint8(writer, v)
	case uint16:
		return encodeUint16(writer, v)
	case uint32:
		return encodeUint32(writer, v)
	case uint64:
		return encodeUint64(writer, v)
	case uint:
		return encodeUint(writer, v)
	case int8:
		return encodeInt8(writer, v)
	case int16:
		return encodeInt16(writer, v)
	case int32:
		return encodeInt32(writer, v)
	case int64:
		return encodeInt64(writer, v)
	case int:
		return encodeInt(writer, v)
	case string:
		return encodeString(writer, v)
	default:
		return encodeValue(writer, reflect.ValueOf(value))
	}

	return 0, nil
}
