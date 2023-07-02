package window

import "fmt"

// Iterator 迭代滑动窗口中的桶信息
type Iterator struct {
	count         int
	iteratedCount int
	cur           *Bucket
}

// Next 判断是否所有桶已被迭代
func (i *Iterator) Next() bool {
	return i.count != i.iteratedCount
}

// Bucket 获取当前桶信息
func (i *Iterator) Bucket() Bucket {
	if !(i.Next()) {
		panic(fmt.Errorf("stat/metric: iteration out of range iteratedCount: %d count: %d", i.iteratedCount, i.count))
	}
	bucket := *i.cur
	i.iteratedCount++
	i.cur = i.cur.Next()
	return bucket
}
