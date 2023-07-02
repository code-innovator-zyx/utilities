package breaker

import (
	"errors"
)

// ErrNotAllowed 熔断器打开
var ErrNotAllowed = errors.New("request not allowed because breaker open")

// Breaker 熔断器实现接口
type Breaker interface {
	Pass() error // 判断熔断器是否允许通过
	// Success 标记调用成功
	Success()
	// Failed 标记调用失败
	Failed()
}
