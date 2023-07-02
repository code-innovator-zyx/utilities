package window

import (
	"time"
)

// Metric 指标接口
type Metric interface {
	// Add 添加计数器
	Add(int64)
	// Value 获取当前值
	Value() int64
}

// Aggregation 常见聚合函数
type Aggregation interface {
	// Min 获取滑动窗口最小值
	Min() float64
	// Max 获取滑动窗口最大值
	Max() float64
	// Avg 计算滑动窗口平均值
	Avg() float64
	// Sum 总数
	Sum() float64
}

// RollingCounter 基于时间段的滑动窗口
// e.g. [[1], [3], [5]]
type RollingCounter interface {
	Metric      // 继承指标接口函数
	Aggregation // 继承聚合函数方法

	Timespan() int
	// Reduce 通过iterator迭代所有窗口进行递减操作
	Reduce(func(Iterator) float64) float64
}

// RollingCounterOpts 滚动计数器
type RollingCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}

type rollingCounter struct {
	policy *RollingPolicy
}

// NewRollingCounter 滚动计数器实例化
func NewRollingCounter(opts RollingCounterOpts) RollingCounter {
	// 新建一个滑动窗口
	window := NewWindow(Options{Size: opts.Size})
	policy := NewRollingPolicy(window, RollingPolicyOpts{BucketDuration: opts.BucketDuration})
	return &rollingCounter{
		policy: policy,
	}
}

func (r *rollingCounter) Add(val int64) {
	r.policy.Add(float64(val))
}

func (r *rollingCounter) Reduce(f func(Iterator) float64) float64 {
	return r.policy.Reduce(f)
}

func (r *rollingCounter) Avg() float64 {
	return r.policy.Reduce(Avg)
}

func (r *rollingCounter) Min() float64 {
	return r.policy.Reduce(Min)
}

func (r *rollingCounter) Max() float64 {
	return r.policy.Reduce(Max)
}

func (r *rollingCounter) Sum() float64 {
	return r.policy.Reduce(Sum)
}

func (r *rollingCounter) Value() int64 {
	return int64(r.Sum())
}

func (r *rollingCounter) Timespan() int {
	r.policy.mu.RLock()
	defer r.policy.mu.RUnlock()
	return r.policy.timespan()
}
