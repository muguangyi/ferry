// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

type txtSerializer struct {
}

func (t *txtSerializer) Marshal(obj interface{}) []byte {
	txt := obj.(string)
	return []byte(txt)
}

func (t *txtSerializer) Unmarshal(data []byte) interface{} {
	return string(data[:])
}

func (t *txtSerializer) Slice(source []byte) int {
	for i, c := range source {
		if '\000' == c {
			return (i + 1)
		}
	}

	return 0
}
