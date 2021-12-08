package uarray

// SlicePage 对长度totalCount的slice进行分页
func SlicePage(page, pageSize int, totalCount int) (sliceStart, sliceEnd int) {
	if page < 0 {
		return 0, 0
	}
	if pageSize <= 0 {
		return 0, 0
	}

	sliceStart = page * pageSize
	sliceEnd = (page + 1) * pageSize

	if sliceStart > totalCount {
		sliceStart = totalCount
	}
	if sliceEnd > totalCount {
		sliceEnd = totalCount
	}

	return sliceStart, sliceEnd
}
