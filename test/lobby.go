// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"sync"

	"github.com/muguangyi/ferry"
)

func newLobby(wg *sync.WaitGroup) ILobby {
	return &lobby{
		wg: wg,
	}
}

type ILobby interface {
}

type lobby struct {
	ferry.Feature
	wg *sync.WaitGroup
}

func (l *lobby) OnInit(s ferry.ISlot) {
	s.Book("IGame")
}

func (l *lobby) OnStart(s ferry.ISlot) {
	s.Visit("IGame").(IGame).Start("level1")
	l.wg.Done()
}
