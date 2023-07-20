package breaker

import (
	"errors"
	"fmt"
)

type State int

const (
	StateClosed   State = iota // 关闭
	StateOpen                  // 开启
	StateHalfOpen              // 半开
)

var (
	ErrStateOpen     = errors.New("circuit breaker is open, drop request")
	ErrStateHalfOpen = errors.New("circuit breaker is half-open, too many calls")
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return fmt.Sprintf("unknown state: %d", s)
	}
}

// IsBreakerError 是否触发熔断错误
func IsBreakerError(err error) bool {
	if errors.Is(err, ErrStateOpen) || errors.Is(err, ErrStateHalfOpen) {
		return true
	}
	return false
}
