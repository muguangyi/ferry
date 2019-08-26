// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec

import (
	"io"
)

// IAny present any value.
type IAny interface {
	// Any returns the source value in IAny object.
	Any() interface{}

	// Bool try to convert the source value to bool type.
	Bool() (bool, error)

	// String try to convert the source value to string type.
	String() (string, error)

	// Int try to convert the source value to int type.
	Int() (int, error)

	// Int8 try to convert the source value to int8 type.
	Int8() (int8, error)

	// Int16 try to convert the source value to int16 type.
	Int16() (int16, error)

	// Int32 try to convert the source value to int32 type.
	Int32() (int32, error)

	// Int64 try to convert the source value to int64 type.
	Int64() (int64, error)

	// Uint try to convert the source value to uint type.
	Uint() (uint, error)

	// Uint8 try to convert the source value to uint8 type.
	Uint8() (uint8, error)

	// Uint16 try to convert the source value to uint16 type.
	Uint16() (uint16, error)

	// Uint32 try to convert the source value to uint32 type.
	Uint32() (uint32, error)

	// Uint64 try to convert the source value to uint64 type.
	Uint64() (uint64, error)

	// Float32 try to convert the source value to float32 type.
	Float32() (float32, error)

	// Float64 try to convert the source value to float64 type.
	Float64() (float64, error)

	// Arr try to convert the source value to []interface{} type.
	Arr() ([]interface{}, error)

	// Map try to convert the source value to map[interface{}]interface{} type.
	Map() (map[interface{}]interface{}, error)

	// Encode IAny object into writer.
	Encode(writer io.Writer) error

	// Decode reader value into IAny object.
	Decode(reader io.Reader) error
}

// NewAny create an IAny object contains value.
func NewAny(value interface{}) IAny {
	a := &any{s: value}
	a.normalize()

	return a
}
