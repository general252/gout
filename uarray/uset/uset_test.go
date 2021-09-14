package uset

import (
	"log"
)

// getResult data中缺少的, data中多余的
func getResult(target *Set, data *Set) (addSet *Set, delSet *Set) {
	addSet = target.Minus(data) // data中缺少的
	delSet = data.Minus(target) // data中多余的
	return
}

func ExampleSet_Minus() {
	var (
		target = New(2, 4, 5, 6, 7) // 最终的结果
		data   = New(2, 3, 4, 5)    // 现有的值, 现有的值如果修改呢
		// 去掉1,2,3,4, 增加6,7
	)

	var c = data.Minus(target) // data中多余的
	var d = target.Minus(data) // data缺少的

	log.Println(c.SortList())
	log.Println(d.SortList())

	a, b := getResult(target, data)
	log.Println(a.SortList(), b.SortList())

	// output:
	//
}
