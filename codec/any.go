// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec

import (
	"fmt"
	"io"
	"reflect"
)

type aType byte

const (
	aNil     aType = 0x00
	aBool    aType = 0x01
	aString  aType = 0x02
	aInt8    aType = 0x03
	aUint8   aType = 0x04
	aInt16   aType = 0x05
	aUint16  aType = 0x06
	aInt32   aType = 0x07
	aUint32  aType = 0x08
	aInt64   aType = 0x0a
	aUint64  aType = 0x0b
	aFloat32 aType = 0x0c
	aFloat64 aType = 0x0d
	aArr     aType = 0x0e
	aMap     aType = 0x0f
)

type any struct {
	s  interface{}
	tp aType
}

func (a *any) Any() interface{} {
	return a.s
}

func (a *any) Bool() (bool, error) {
	if aBool == a.tp {
		return a.s.(bool), nil
	}

	return false, fmt.Errorf("Can't convert %d to bool!", a.tp)
}

func (a *any) String() (string, error) {
	if aString == a.tp {
		return a.s.(string), nil
	}

	return "", fmt.Errorf("Can't convert %d to string!", a.tp)
}

func (a *any) Int() (int, error) {
	switch a.tp {
	case aInt8:
		return int(a.s.(int8)), nil
	case aInt16:
		return int(a.s.(int16)), nil
	case aInt32:
		return int(a.s.(int32)), nil
	case aInt64:
		return int(a.s.(int64)), nil
	}

	return -1, fmt.Errorf("Can't convert %d to int!", a.tp)
}

func (a *any) Int8() (int8, error) {
	switch a.tp {
	case aInt8:
		return int8(a.s.(int8)), nil
	case aInt16:
		return int8(a.s.(int16)), nil
	case aInt32:
		return int8(a.s.(int32)), nil
	case aInt64:
		return int8(a.s.(int64)), nil
	}

	return -1, fmt.Errorf("Can't convert %d to int8!", a.tp)
}

func (a *any) Int16() (int16, error) {
	switch a.tp {
	case aInt8:
		return int16(a.s.(int8)), nil
	case aInt16:
		return int16(a.s.(int16)), nil
	case aInt32:
		return int16(a.s.(int32)), nil
	case aInt64:
		return int16(a.s.(int64)), nil
	}

	return -1, fmt.Errorf("Can't convert %d to int16!", a.tp)
}

func (a *any) Int32() (int32, error) {
	switch a.tp {
	case aInt8:
		return int32(a.s.(int8)), nil
	case aInt16:
		return int32(a.s.(int16)), nil
	case aInt32:
		return int32(a.s.(int32)), nil
	case aInt64:
		return int32(a.s.(int64)), nil
	}

	return -1, fmt.Errorf("Can't convert %d to int32!", a.tp)
}

func (a *any) Int64() (int64, error) {
	switch a.tp {
	case aInt8:
		return int64(a.s.(int8)), nil
	case aInt16:
		return int64(a.s.(int16)), nil
	case aInt32:
		return int64(a.s.(int32)), nil
	case aInt64:
		return int64(a.s.(int64)), nil
	}

	return -1, fmt.Errorf("Can't convert %d to int64!", a.tp)
}

func (a *any) Uint() (uint, error) {
	switch a.tp {
	case aUint8:
		return uint(a.s.(uint8)), nil
	case aUint16:
		return uint(a.s.(uint16)), nil
	case aUint32:
		return uint(a.s.(uint32)), nil
	case aUint64:
		return uint(a.s.(uint64)), nil
	}

	return 0, fmt.Errorf("Can't convert %d to uint8!", a.tp)
}

func (a *any) Uint8() (uint8, error) {
	switch a.tp {
	case aUint8:
		return uint8(a.s.(uint8)), nil
	case aUint16:
		return uint8(a.s.(uint16)), nil
	case aUint32:
		return uint8(a.s.(uint32)), nil
	case aUint64:
		return uint8(a.s.(uint64)), nil
	}

	return 0, fmt.Errorf("Can't convert %d to uint8!", a.tp)
}

func (a *any) Uint16() (uint16, error) {
	switch a.tp {
	case aUint8:
		return uint16(a.s.(uint8)), nil
	case aUint16:
		return uint16(a.s.(uint16)), nil
	case aUint32:
		return uint16(a.s.(uint32)), nil
	case aUint64:
		return uint16(a.s.(uint64)), nil
	}

	return 0, fmt.Errorf("Can't convert %d to uint16!", a.tp)
}

func (a *any) Uint32() (uint32, error) {
	switch a.tp {
	case aUint8:
		return uint32(a.s.(uint8)), nil
	case aUint16:
		return uint32(a.s.(uint16)), nil
	case aUint32:
		return uint32(a.s.(uint32)), nil
	case aUint64:
		return uint32(a.s.(uint64)), nil
	}

	return 0, fmt.Errorf("Can't convert %d to uint32!", a.tp)
}

func (a *any) Uint64() (uint64, error) {
	switch a.tp {
	case aUint8:
		return uint64(a.s.(uint8)), nil
	case aUint16:
		return uint64(a.s.(uint16)), nil
	case aUint32:
		return uint64(a.s.(uint32)), nil
	case aUint64:
		return uint64(a.s.(uint64)), nil
	}

	return 0, fmt.Errorf("Can't convert %d to uint64!", a.tp)
}

func (a *any) Float32() (float32, error) {
	switch a.tp {
	case aFloat32:
		return float32(a.s.(float32)), nil
	case aFloat64:
		return float32(a.s.(float64)), nil
	}

	return -1, fmt.Errorf("Can't convert %d to float32!", a.tp)
}

func (a *any) Float64() (float64, error) {
	switch a.tp {
	case aFloat32:
		return float64(a.s.(float32)), nil
	case aFloat64:
		return float64(a.s.(float64)), nil
	}

	return -1, fmt.Errorf("Can't convert %d to float64!", a.tp)
}

func (a *any) Arr() ([]interface{}, error) {
	if aArr == a.tp {
		return a.s.([]interface{}), nil
	}

	return nil, fmt.Errorf("Can't convert %d to array!", a.tp)
}

func (a *any) Map() (map[interface{}]interface{}, error) {
	if aMap == a.tp {
		return a.s.(map[interface{}]interface{}), nil
	}

	return nil, fmt.Errorf("Can't convert %d to map!", a.tp)
}

func (a *any) Encode(writer io.Writer) error {
	_, err := encode(writer, a.s)
	return err
}

func (a *any) Decode(reader io.Reader) error {
	var err error
	a.s, err = decode(reader)
	if err != nil {
		return err
	}

	a.normalize()

	return nil
}

func (a *any) normalize() {
	if nil == a.s {
		a.tp = aNil
	}

	switch a.s.(type) {
	case bool:
		a.tp = aBool
	case string:
		a.tp = aString
	case int8:
		a.tp = aInt8
	case uint8:
		a.tp = aUint8
	case int16:
		a.tp = aInt16
	case uint16:
		a.tp = aUint16
	case int32:
		a.tp = aInt32
	case uint32:
		a.tp = aUint32
	case int64:
		a.tp = aInt64
	case uint64:
		a.tp = aUint64
	case float32:
		a.tp = aFloat32
	case float64:
		a.tp = aFloat64
	default:
		vt := reflect.ValueOf(a.s)
		switch vt.Kind() {
		case reflect.Array, reflect.Slice:
			a.tp = aArr
		case reflect.Map:
			a.tp = aMap
		}
	}
}
