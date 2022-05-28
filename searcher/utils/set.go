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

func DeleteDuplicatedWord(words []string) []string {
	wordSet := make(Set)
	var res []string = make([]string, 0, len(words)/2)
	for _, word := range words {
		if wordSet.IsExist(word) || word == "" {
			continue
		}

		res = append(res, word)
	}

	return res
}
