// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chancall

import (
	"reflect"
)

func newMeta(name string, target interface{}) *meta {
	m := new(meta)
	m.name = name
	m.funcs = make(map[string]*fcall)
	m.collect(target)

	return m
}

type meta struct {
	name  string
	funcs map[string]*fcall
}

type fcall struct {
	fn      *reflect.Value
	timeout float32
}

func (m *meta) collect(target interface{}) {
	value := reflect.ValueOf(target)
	t := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		fn := value.Method(i)
		m.funcs[t.Method(i).Name] = &fcall{
			fn:      &fn,
			timeout: cDefaultTimeout,
		}
	}
}

func (m *meta) call(method string, args ...interface{}) []interface{} {
	f := m.funcs[method]
	if nil != f && f.fn.IsValid() {
		params := make([]reflect.Value, 0)
		for _, arg := range args {
			params = append(params, reflect.ValueOf(arg))
		}
		ret := f.fn.Call(params)

		result := make([]interface{}, len(ret))
		for i, r := range ret {
			result[i] = r.Interface()
		}

		return result
	}

	return nil
}

func (m *meta) timeout(method string) float32 {
	f := m.funcs[method]
	if nil != f {
		return f.timeout
	}

	return cDefaultTimeout
}

func (m *meta) setTimeout(method string, timeout float32) {
	f := m.funcs[method]
	if nil != f {
		f.timeout = timeout
	}
}

const (
	cDefaultTimeout float32 = 1.0
)
