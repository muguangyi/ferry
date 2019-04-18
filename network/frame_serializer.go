// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package network

type frameSerializer struct {
}

func (f *frameSerializer) Marshal(obj interface{}) []byte {
	return nil
}

func (f *frameSerializer) Unmarshal(data []byte) interface{} {
	return nil
}

func (f *frameSerializer) Slice(source []byte) int {
	return 0
}
