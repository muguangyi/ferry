// Copyright 2019 MuGuangyi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package misc

// ISet interface
type ISet interface {
	// Add item to set.
	Add(item interface{}) bool

	// Remove item from set and return success or not.
	Remove(item interface{}) bool

	// Convert to a slice.
	ToSlice() []interface{}
}

// Create an empty set.
func NewSet() ISet {
	s := new(set)
	s.data = make(map[interface{}]bool)

	return s
}

type set struct {
	data map[interface{}]bool
}

func (s *set) Add(item interface{}) bool {
	if !s.data[item] {
		s.data[item] = true
		return true
	}

	return false
}

func (s *set) Remove(item interface{}) bool {
	if s.data[item] {
		delete(s.data, item)
		return true
	}

	return false
}

func (s *set) ToSlice() []interface{} {
	arr := make([]interface{}, 0)
	for k := range s.data {
		arr = append(arr, k)
	}

	return arr
}
