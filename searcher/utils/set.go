package utils

// 借助Map自定义set结构

type Set map[interface{}]struct{}

func (s Set) Add(key interface{}) {
	s[key] = struct{}{}
}

func (s Set) IsExist(key interface{}) bool {
	if _, ok := s[key]; !ok {
		return false
	}

	return true
}

func (s Set) DeleteKey(key interface{}) {
	delete(s, key)
}
