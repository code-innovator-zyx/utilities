package sre

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"utilities/breaker"
	"utilities/window"
)

// Option 熔断器配置
type Option func(*options)

const (
	// 断路器打开时的状态打开，睡眠后不允许请求

	StateOpen int32 = iota
	// 断路器打开时的状态关，睡眠后不允许请求

	StateClosed
)

// options 熔断器配置参数.
type options struct {
	k       float64
	request int64
	bucket  int
	window  time.Duration
}

/*
	SetK  K  默认2
	降低 K 值会使自适应限流算法更加激进（允许客户端在算法启动时拒绝更多本地请求）
	增加 K 值会使自适应限流算法不再那么激进（允许服务端在算法启动时尝试接收更多的请求，与上面相反）
*/
func SetK(k float64) Option {
	return func(c *options) {
		c.k = k
	}
}

// SetRequest 允许的最小请求数
func SetRequest(r int64) Option {
	return func(c *options) {
		c.request = r
	}
}

// SetWindow 滑动窗口的大小
func SetWindow(d time.Duration) Option {
	return func(c *options) {
		c.window = d
	}
}

// SetBucket 对时间范围内的滑动窗口桶设置编号
func SetBucket(b int) Option {
	return func(c *options) {
		c.bucket = b
	}
}

// Breaker google sre 算法实例   需要实现熔断器三个函数
type Breaker struct {
	stat window.RollingCounter // 继承了滑动窗口RollingCounter的所有函数方法 主要是进行 滑动窗口，记录成功失败
	r    *rand.Rand            // 返回一个取值范围在 [0.0, 1.0) 的伪随机 Float64 值
	// rand.New(...) returns a non thread safe object
	randLock sync.Mutex
	/*
		降低 K 值会使自适应限流算法更加激进（允许客户端在算法启动时拒绝更多本地请求）
		增加 K 值会使自适应限流算法不再那么激进（允许服务端在算法启动时尝试接收更多的请求，与上面相反）
	*/
	k       float64 // 成功系数  requests = success * k
	request int64   // 触发熔断的最小请求数  当总数 < request时，不判断是否熔断

	state int32 // 熔断器状态 打开或者关闭
}

// BreakerGroup 原子引擎熔断需要针对每一个原子引擎，需要一对一的breaker 熔断器
type BreakerGroup struct {
	opts []Option
	new  func(...Option) breaker.Breaker
	objs map[string]breaker.Breaker
	sync.RWMutex
}

// NewBreaker 实例化breaker 对象
func NewBreaker(opts ...Option) breaker.Breaker {
	opt := options{
		k:       2,               // 默认2
		request: 100,             // 触发熔断的最少请求数量（请求少于该值时不会触发熔断）
		bucket:  10,              // 统计桶大小
		window:  5 * time.Second, // 统计桶窗口时间
	}
	for _, o := range opts {
		o(&opt)
	}
	// 滚动计数器
	stat := window.NewRollingCounter(window.RollingCounterOpts{
		Size:           opt.bucket,
		BucketDuration: time.Duration(int64(opt.window) / int64(opt.bucket)),
	})
	return &Breaker{
		stat:    stat,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		request: opt.request,
		k:       opt.k,
		state:   StateClosed, // 初始关闭状态
	}
}

// NewGroup 获取熔断器组
func NewGroup(opts ...Option) *BreakerGroup {
	return &BreakerGroup{
		opts: opts,
		new:  NewBreaker,
		objs: make(map[string]breaker.Breaker),
	}
}

// 累加滑动窗口 Bucket 中的 bucket.Count 和 bucket.Points：
func (b *Breaker) stats() (success int64, total int64) {
	b.stat.Reduce(func(iterator window.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			total += bucket.Count
			for _, p := range bucket.Points {
				success += int64(p)
			}
		}
		return 0
	})
	return
}

/**
	算法详情:
在通常情况下（无错误发生时） requests==accepts
当后端出现异常情况时，accepts 的数量会逐渐小于 requests
当后端持续异常时，客户端可以继续发送请求直到 requests=K∗accepts，一旦超过这个值，客户端就启动自适应限流机制，新产生的请求在本地会以 p 概率（下面描述的 Client request rejection probability 定义）被拒绝
当客户端主动丢弃请求时，requests 值会一直增大，在某个时间点会超过 K∗accepts，使概率 p 计算出来的值大于 0，此时客户端会以此概率对请求做主动丢弃
当后端逐渐恢复时，accepts 增加，（同时 requests 值也会增加，但是由于 K 的关系，K×accepts 的放大倍数更快），使得 requests−K×acceptsrequests+1 变为负数，从而概率 p==0，客户端自适应限流结束
*/

// Pass 判断请求是否允许通过   允许返回nil
func (b *Breaker) Pass() error {
	// 统计滑动窗口时间段内调用成功量和总量
	accepts, total := b.stats()
	requests := b.k * float64(accepts)
	// 如果总量小于请求量 或者客户端可以发送的量
	if total < b.request || float64(total) < requests {
		// 设置关闭状态
		atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClosed)
		return nil
	}
	// 否则熔断器状态开启
	atomic.CompareAndSwapInt32(&b.state, StateClosed, StateOpen)
	// 通过 google sre breaker 算法获取丢弃率
	drop := math.Max(0, (float64(total)-requests)/float64(total+1))
	if b.discarded(drop) {
		return breaker.ErrNotAllowed
	}
	return nil
}

// Success 标记调用成功
func (b *Breaker) Success() {
	b.stat.Add(1)
}

// Failed 标记调用失败
func (b *Breaker) Failed() {
	// 当调用失败失败时，继续添加失败数量，让丢弃率更高
	b.stat.Add(0)
}

// 判断是否忽略当前请求
func (b *Breaker) discarded(probably float64) (truth bool) {
	b.randLock.Lock()
	defer b.randLock.Unlock()
	// probably 越大时， 忽略的概率就越大。
	return b.r.Float64() < probably
}

// Get 通过原子引擎名称获取一个熔断器  如果已经存在  则返回原有的  不存则创建一个新的
func (g *BreakerGroup) Get(key string) breaker.Breaker {
	g.RLock()
	obj, ok := g.objs[key]
	if ok {
		g.RUnlock()
		return obj
	}
	g.RUnlock()

	// double check
	g.Lock()
	defer g.Unlock()
	obj, ok = g.objs[key]
	if ok {
		return obj
	}
	obj = g.new(g.opts...)
	g.objs[key] = obj
	return obj
}

// Clear 可以作为后续手动关闭限流器的控制操作
func (g *BreakerGroup) Clear() {
	g.Lock()
	g.objs = make(map[string]breaker.Breaker)
	g.Unlock()
}
