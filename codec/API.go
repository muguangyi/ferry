// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package codec

import (
	"io"
)

// Encode target value to stream writer.
func Encode(writer io.Writer, value interface{}) (int, error) {
	return encode(writer, value)
}

// Decode value from stream reader.
func Decode(reader io.Reader) (interface{}, error) {
	return decode(reader)
}
