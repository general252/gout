package uarray

// []string deduplicate
func StringsDeduplicate(array []string) []string {
	var arr = make([]string, 0)
	var m = make(map[string]bool)
	for _, d := range array {
		_, ok := m[d]
		if !ok {
			m[d] = true
			arr = append(arr, d)
		}
	}
	return arr
}

// []int deduplicate
func Int32sDeduplicate(array []int) []int {
	var arr = make([]int, 0)
	var m = make(map[int]bool)
	for _, d := range array {
		_, ok := m[d]
		if !ok {
			m[d] = true
			arr = append(arr, d)
		}
	}
	return arr
}

// []int64 deduplicate
func Int64sDeduplicate(array []int64) []int64 {
	var arr = make([]int64, 0)
	var m = make(map[int64]bool)
	for _, d := range array {
		_, ok := m[d]
		if !ok {
			m[d] = true
			arr = append(arr, d)
		}
	}
	return arr
}
