// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/muguangyi/seek/seek"
)

type IMath interface {
	Add(x float64, y float64) float64
}

func newMath() IMath {
	return &math{}
}

type math struct {
	seek.Signal
}

func (math *math) Print(msg string) {
	time.Sleep(2 * time.Second)
	fmt.Println("print...")
	fmt.Println(msg)
}

func (math *math) Add(x float64, y float64) float64 {
	fmt.Println("add...")
	result := x + y
	return result
}
