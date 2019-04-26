// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"sync"

	"github.com/muguangyi/unite/unite"
)

func newLobby(wg *sync.WaitGroup) ILobby {
	return &LobbyControl{
		wg: wg,
	}
}

type ILobby interface {
}

type LobbyControl struct {
	unite.UnitControl
	wg *sync.WaitGroup
}

func (l *LobbyControl) OnInit(u unite.IUnit) {
	l.UnitControl.OnInit(u)
	l.Import("IGame")
}

func (l *LobbyControl) OnStart() {
	l.Call("IGame", "Start", "level1")
	l.wg.Done()
}
