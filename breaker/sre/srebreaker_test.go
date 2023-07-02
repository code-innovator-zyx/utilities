package sre

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
	"utilities/breaker"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2022/11/14 18:27
* @Package:
 */

func Test_CircuitBreaker(t *testing.T) {
	// 单例模式
	t.Run("test circuit", func(t *testing.T) {
		b := NewBreaker()
		go func(bk breaker.Breaker) {
			for i := 0; i < 100; i++ {
				//标记成功
				bk.Success()
			}
		}(b)

		go func(bk breaker.Breaker) {
			for i := 0; i < 300; i++ {
				bk.Failed()
			}
		}(b)
		time.Sleep(time.Second * 1)
		err := b.Pass()
		fmt.Printf("err=%v", err)
	})

	t.Run("test compare", func(t *testing.T) {
		var stat int32 = 0
		atomic.CompareAndSwapInt32(&stat, StateClosed, StateOpen)
		fmt.Println(stat)
	})
	// 熔断器组模式
	t.Run("test group", func(t *testing.T) {
		bg := NewGroup(SetK(2), SetRequest(20))
		b := bg.Get("image_od_vip")
		go func(bk breaker.Breaker) {
			for i := 0; i < 100; i++ {
				//标记成功
				bk.Success()
			}
		}(b)

		go func(bk breaker.Breaker) {
			for i := 0; i < 300; i++ {
				bk.Failed()
			}
		}(b)
		time.Sleep(time.Second * 1)
		err := b.Pass()
		fmt.Printf("err=%v", err)
	})
}
