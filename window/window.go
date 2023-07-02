package window

// Bucket 桶信息
type Bucket struct {
	Points []float64
	Count  int64   // 计算
	next   *Bucket // 指向下一个记录的指针
}

// Append 向桶内追加数据
func (b *Bucket) Append(val float64) {
	b.Points = append(b.Points, val)
	b.Count++
}

// Add 向points 添加偏移量
func (b *Bucket) Add(offset int, val float64) {
	b.Points[offset] += val
	b.Count++
}

// Reset 重置桶数据
func (b *Bucket) Reset() {
	b.Points = b.Points[:0]
	b.Count = 0
}

// Next 获取指向下一个桶的指针地址
func (b *Bucket) Next() *Bucket {
	return b.next
}

// Window 滑动窗口包含多个请求桶
type Window struct {
	buckets []Bucket
	size    int
}

// Options 创建滑动窗口的参数
type Options struct {
	Size int
}

// NewWindow creates a new Window
func NewWindow(opts Options) *Window {
	buckets := make([]Bucket, opts.Size)
	for offset := range buckets {
		buckets[offset].Points = make([]float64, 0)
		nextOffset := offset + 1
		if nextOffset == opts.Size {
			nextOffset = 0
		}
		buckets[offset].next = &buckets[nextOffset]
	}
	return &Window{buckets: buckets, size: opts.Size}
}

// ResetWindow 重置滑动窗口
func (w *Window) ResetWindow() {
	for offset := range w.buckets {
		w.ResetBucket(offset)
	}
}

// ResetBucket 重置指定偏移
func (w *Window) ResetBucket(offset int) {
	w.buckets[offset%w.size].Reset()
}

// ResetBuckets 重置存储桶
func (w *Window) ResetBuckets(offset int, count int) {
	for i := 0; i < count; i++ {
		w.ResetBucket(offset + i)
	}
}

// Append 向窗口中指定的桶追加
func (w *Window) Append(offset int, val float64) {
	w.buckets[offset%w.size].Append(val)
}

// Add 将给定值添加到存储桶中索引等于给定偏移量的最新点。
func (w *Window) Add(offset int, val float64) {
	offset %= w.size
	if w.buckets[offset].Count == 0 {
		w.buckets[offset].Append(val)
		return
	}
	w.buckets[offset].Add(0, val)
}

// Bucket 通过偏移量获取bucket
func (w *Window) Bucket(offset int) Bucket {
	return w.buckets[offset%w.size]
}

// Size 获取窗口大小
func (w *Window) Size() int {
	return w.size
}

// Iterator 通过偏移量返回计数桶数迭代器。
func (w *Window) Iterator(offset int, count int) Iterator {
	return Iterator{
		count: count,
		cur:   &w.buckets[offset%w.size],
	}
}
