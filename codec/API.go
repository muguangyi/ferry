// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec

import (
	"io"
)

// IAny present any value.
type IAny interface {
	Any() interface{}
	Bool() (bool, error)
	String() (string, error)
	Int() (int, error)
	Int8() (int8, error)
	Int16() (int16, error)
	Int32() (int32, error)
	Int64() (int64, error)
	Uint() (uint, error)
	Uint8() (uint8, error)
	Uint16() (uint16, error)
	Uint32() (uint32, error)
	Uint64() (uint64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Arr() ([]interface{}, error)
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
