package main

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func maxIntSlice(a []int) int {
	m := 0
	for _, x := range a {
		if x > m {
			m = x
		}
	}
	return m
}

// Sorts values by keys
type sortByKey struct {
	key   []int
	value []int
}

func (s sortByKey) Len() int {
	return len(s.value)
}
func (s sortByKey) Swap(i, j int) {
	s.key[i], s.key[j] = s.key[j], s.key[i]
	s.value[i], s.value[j] = s.value[j], s.value[i]
}
func (s sortByKey) Less(i, j int) bool {
	return s.key[i] < s.key[j]
}

func NewSortByKey(key, value []int) sortByKey {
	keyCopy := make([]int, len(key))
	copy(keyCopy, key)
	valueCopy := make([]int, len(value))
	copy(valueCopy, value)

	return sortByKey{keyCopy, valueCopy}
}
