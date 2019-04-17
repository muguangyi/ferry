// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package unit

type IUnit interface {
	OnInit()
	OnDestroy()
	OnUpdate(closeSig chan bool)
}
