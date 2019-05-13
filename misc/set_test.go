// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package misc_test

import (
	"testing"

	"github.com/muguangyi/ferry/misc"
)

func TestAddAndRemove(t *testing.T) {
	set := misc.NewSet()
	if !set.Add(1) {
		t.Fail()
	}

	if set.Add(1) {
		t.Fail()
	}

	if !set.Remove(1) {
		t.Fail()
	}
}

func TestToSlice(t *testing.T) {
	set := misc.NewSet()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	slice := set.ToSlice()
	if 3 != len(slice) {
		t.Fail()
	}
}
