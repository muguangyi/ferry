// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
)

func main() {
	flag.Parse()

	ioutil.WriteFile("h.gen.go", bytes.NewBufferString("yes").Bytes(), 0644)
}
