package breaker

import (
	"sync"
	"time"
)

// Breaker 熔断器
type Breaker struct {
	name            string        // 熔断器名称
	state           State         // 熔断器状态
	halfOpenMaxCall uint64        // 半开期间最大请求数（半开期间，若请求前的总请求数大于此则丢弃，若请求后的连续成功数大于此则关闭熔断器）
	mu              sync.RWMutex  // 互斥锁
	openTime        time.Time     // 熔断器打开时间
	windowInterval  time.Duration // 窗口间隔
	coolDownTime    time.Duration // 冷却时间（从开到半开的时间间隔）
	metric          Metric        // 指标
	strategyFn      StrategyFn    // 熔断策略
}

const (
	DefaultWindowInterval          = time.Second // 默认窗口间隔
	DefaultCoolDownTime            = time.Second // 默认冷却时间
	DefaultHalfOpenMaxCall         = 5           // 默认半开期间最大请求数
	DefaultFailThreshold           = 10          // 默认失败数阈值
	DefaultContinuousFailThreshold = 10          // 默认连续失败数阈值
	DefaultFailRate                = 0.6         // 默认失败率阈值
	DefaultMinCall                 = 10          // 默认失败率策略的最小请求数
)

var defaultBreaker = Breaker{
	windowInterval:  DefaultWindowInterval,
	coolDownTime:    DefaultCoolDownTime,
	halfOpenMaxCall: DefaultHalfOpenMaxCall,
	strategyFn:      FailStrategyFn(DefaultFailThreshold),
}

// NewBreaker returns a Breaker object.
// opts can be used to customize the Breaker.
func NewBreaker(opts ...Option) *Breaker {
	breaker := &defaultBreaker
	for _, opt := range opts {
		opt(breaker)
	}
	if breaker.name == "" {
		breaker.name = "breakerName"
	}
	breaker.newWindow(time.Now())
	return breaker
}

// Do do fn
func (b *Breaker) Do(fn func() error) error {
	// before call
	batch, err := b.beforeCall()
	if err != nil {
		return err
	}

	// recover
	defer func() {
		if err := recover(); err != nil {
			b.afterCall(batch, false)
			//panic(err)
		}
	}()

	// call function
	err = fn()

	// after call
	b.afterCall(batch, err == nil)

	return err
}

// beforeCall before call
func (b *Breaker) beforeCall() (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.updateState(time.Now())

	switch b.state {
	case StateOpen:
		return b.metric.WindowBatch, ErrStateOpen
	case StateHalfOpen:
		// 请求数 ≥ 半开最大请求数，丢弃请求
		if b.metric.TotalRequest >= b.halfOpenMaxCall {
			return b.metric.WindowBatch, ErrStateHalfOpen
		}
	}

	b.metric.onRequest()
	return b.metric.WindowBatch, nil
}

// after call
func (b *Breaker) afterCall(batch uint64, result bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	b.updateState(now)
	if b.metric.WindowBatch != batch {
		return
	}

	if result {
		b.onSuccess(now)
	} else {
		b.onFail(now)
	}
}

// newWindow create new window
func (b *Breaker) newWindow(t time.Time) {
	b.metric.NewWindowBatch()
	b.metric.OnReset()

	switch b.state {
	case StateClosed:
		if b.windowInterval <= 0 {
			b.metric.WindowExpiry = time.Time{}
		} else {
			b.metric.WindowExpiry = t.Add(b.windowInterval)
		}
	case StateOpen:
		b.metric.WindowExpiry = t.Add(b.coolDownTime)
	default:
		b.metric.WindowExpiry = time.Time{}
	}
}

// onSuccess call on success
func (b *Breaker) onSuccess(t time.Time) {
	b.metric.onSuccess()
	switch b.state {
	case StateClosed:
		if b.metric.WindowExpiry.Before(t) {
			b.newWindow(t)
		}
	case StateHalfOpen:
		if b.metric.ContinuousSuccess >= b.halfOpenMaxCall {
			b.setState(StateClosed, t)
		}
	}
}

// onFail call on failure
func (b *Breaker) onFail(t time.Time) {
	b.metric.onFail()
	switch b.state {
	case StateClosed:
		if b.strategyFn(b.metric) {
			b.setState(StateOpen, t)
		}
	case StateHalfOpen:
		b.setState(StateOpen, t)
	}
}

// setState set breaker state
func (b *Breaker) setState(state State, t time.Time) {
	if b.state == state {
		return
	}

	b.state = state
	b.openTime = t
	b.newWindow(t)
}

// updateState update state
func (b *Breaker) updateState(t time.Time) {
	switch b.state {
	case StateClosed:
		if !b.metric.WindowExpiry.IsZero() && b.metric.WindowExpiry.Before(t) {
			b.newWindow(t)
		}
	case StateOpen:
		if b.metric.WindowExpiry.Before(t) {
			b.newWindow(t)
			b.setState(StateHalfOpen, t)
		}
	}
}
