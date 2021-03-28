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
