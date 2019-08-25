// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ferry

import (
	"io"
)

// IProto presents proto object.
type IProto interface {
	// Marshal proto object into writer.
	Marshal(writer io.Writer) error

	// Unmarshal reader data into proto object.
	Unmarshal(reader io.Reader) error
}
