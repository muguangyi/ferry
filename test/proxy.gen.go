//
// This code was generated by seek tool.
//
// Changes to this file may cause incorrect behavior and will be lost if the code is regenerated.
//
// 2019-05-10 10:23:10
//

package main

import (
	"github.com/muguangyi/seek"
)

// IGame from: game.go
type igameproxy struct {
	sandbox seek.ISandbox
}

func (p *igameproxy) Start(level string) {
	err := p.sandbox.Call("IGame", "Start", level)

	if nil != err {
		return
	}
}

// end IGame

// ILobby from: lobby.go
type ilobbyproxy struct {
	sandbox seek.ISandbox
}

// end ILobby

// IMath from: math.go
type imathproxy struct {
	sandbox seek.ISandbox
}

func (p *imathproxy) Add(x float64, y float64) float64 {
	results, err := p.sandbox.CallWithResult("IMath", "Add", x, y)

	if nil != err {
		return 0
	}

	return results[0].(float64)
}

func (p *imathproxy) Print(msg string) {
	err := p.sandbox.Call("IMath", "Print", msg)

	if nil != err {
		return
	}
}

// end IMath

// ILogin from: sub\login.go
type iloginproxy struct {
	sandbox seek.ISandbox
}

func (p *iloginproxy) Login(name string, pwd string) bool {
	results, err := p.sandbox.CallWithResult("ILogin", "Login", name, pwd)

	if nil != err {
		return false
	}

	return results[0].(bool)
}

func (p *iloginproxy) Logout() {
	err := p.sandbox.Call("ILogin", "Logout")

	if nil != err {
		return
	}
}

// end ILogin

// Register to seek
var (
	igameproxysucc  bool = seek.Register("IGame", func(sandbox seek.ISandbox) interface{} { return &igameproxy{sandbox: sandbox} })
	ilobbyproxysucc bool = seek.Register("ILobby", func(sandbox seek.ISandbox) interface{} { return &ilobbyproxy{sandbox: sandbox} })
	imathproxysucc  bool = seek.Register("IMath", func(sandbox seek.ISandbox) interface{} { return &imathproxy{sandbox: sandbox} })
	iloginproxysucc bool = seek.Register("ILogin", func(sandbox seek.ISandbox) interface{} { return &iloginproxy{sandbox: sandbox} })
)
