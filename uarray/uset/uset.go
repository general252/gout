package uset

import (
	"sort"
	"sync"
)

type (
	DogData int64
)

func DogDataCompare(a, b DogData) int {
	if a > b {
		return 1
	} else if a == b {
		return 0
	} else {
		return -1
	}
}

type Set struct {
	sync.RWMutex
	m map[DogData]bool
}

// New 新建集合对象
func New(items ...DogData) *Set {
	s := &Set{
		m: make(map[DogData]bool, len(items)),
	}
	s.Add(items...)
	return s
}

// Add 添加元素
func (s *Set) Add(items ...DogData) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// Remove 删除元素
func (s *Set) Remove(items ...DogData) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

// Has 判断元素是否存在
func (s *Set) Has(items ...DogData) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// Count 元素个数
func (s *Set) Count() int {
	return len(s.m)
}

// Clear 清空集合
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[DogData]bool{}
}

// Empty 空集合判断
func (s *Set) Empty() bool {
	return len(s.m) == 0
}

// List 无序列表
func (s *Set) List() []DogData {
	s.RLock()
	defer s.RUnlock()
	list := make([]DogData, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// SortList 排序列表
func (s *Set) SortList() []DogData {
	s.RLock()
	defer s.RUnlock()
	var list []DogData
	for item := range s.m {
		list = append(list, item)
	}

	sort.Slice(list, func(i, j int) bool {
		return DogDataCompare(list[i], list[j]) < 0 // return list[i] < list[j]
	})
	return list
}

// Union 并集
func (s *Set) Union(sets ...*Set) *Set {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range set.m {
			r.m[e] = true
		}
	}
	return r
}

// Minus 差集
func (s *Set) Minus(sets ...*Set) *Set {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range set.m {
			if _, ok := s.m[e]; ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Intersect 交集
func (s *Set) Intersect(sets ...*Set) *Set {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range s.m {
			if _, ok := set.m[e]; !ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Complement 补集
func (s *Set) Complement(full *Set) *Set {
	r := New()
	for e := range full.m {
		if _, ok := s.m[e]; !ok {
			r.Add(e)
		}
	}
	return r
}
