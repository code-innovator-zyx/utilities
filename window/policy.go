package window

import (
	"sync"
	"time"
)

/**
RollingPolicy 基于持续时间段内的滑动窗口策略
*/

type RollingPolicy struct {
	mu     sync.RWMutex // *需要读写锁
	size   int
	window *Window
	offset int

	bucketDuration time.Duration // 时长
	lastAppendTime time.Time     //最后一次的时间
}

// RollingPolicyOpts 创建策略配置参数  待定---后续可追加
type RollingPolicyOpts struct {
	BucketDuration time.Duration
}

// NewRollingPolicy 根据给定的窗口和策略配置选择创建一个新的滚动策略。
func NewRollingPolicy(window *Window, opts RollingPolicyOpts) *RollingPolicy {
	return &RollingPolicy{
		window: window,
		size:   window.Size(),
		offset: 0,

		bucketDuration: opts.BucketDuration,
		lastAppendTime: time.Now(),
	}
}

/**
timespan 返回自上次追加时间以来传递的存储桶编号
如果它比上次记录的时间早一个存储桶持续时间，返回大小。
*/

func (r *RollingPolicy) timespan() int {
	v := int(time.Since(r.lastAppendTime) / r.bucketDuration)
	if v > -1 { // 如果最后的时间小于当前时间 (一般不会出现，除非有Bug,但是为了健壮性，不得不写 ლ(′◉❥◉｀ლ))
		return v
	}
	return r.size
}

//对滑动窗口的当前桶进行操作
func (r *RollingPolicy) apply(f func(offset int, val float64), val float64) {
	// 开启读写锁
	r.mu.Lock()
	defer r.mu.Unlock()
	// 获取当前偏移量
	timespan := r.timespan()
	oriTimespan := timespan
	if timespan > 0 {
		start := (r.offset + 1) % r.size
		end := (r.offset + timespan) % r.size
		if timespan > r.size {
			timespan = r.size
		}
		// 重置过期的桶
		r.window.ResetBuckets(start, timespan)
		r.offset = end
		r.lastAppendTime = r.lastAppendTime.Add(time.Duration(oriTimespan * int(r.bucketDuration)))
	}
	f(r.offset, val)
}

func (r *RollingPolicy) Append(val float64) {
	r.apply(r.window.Append, val)
}

func (r *RollingPolicy) Add(val float64) {
	r.apply(r.window.Add, val)
}

func (r *RollingPolicy) Reduce(f func(Iterator) float64) (val float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	timespan := r.timespan()
	if count := r.size - timespan; count > 0 {
		offset := r.offset + timespan + 1
		if offset >= r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	return val
}
