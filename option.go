package breaker

import (
	"github.com/bluesaka/go-breaker/notify"
	"time"
)

type Option func(o *Breaker)

// StrategyOption 策略选项
type StrategyOption struct {
	Strategy                int     // 策略类型
	FailThreshold           uint64  // 失败数阈值 (失败数策略)
	ContinuousFailThreshold uint64  // 连续失败数阈值 (连续失败数策略)
	FailRate                float64 // 失败率阈值 (失败率策略)
	MinCall                 uint64  // 最小请求数 (失败率策略)
}

// WithName returns a function to set the name of Breaker
func WithName(s string) Option {
	return func(options *Breaker) {
		options.name = s
	}
}

// WithWindowInterval returns a function to set the windowInterval of Breaker
func WithWindowInterval(d time.Duration) Option {
	return func(options *Breaker) {
		options.windowInterval = d
	}
}

// WithCoolDownTime returns a function to set the coolDownTime of Breaker
func WithCoolDownTime(d time.Duration) Option {
	return func(options *Breaker) {
		options.coolDownTime = d
	}
}

// WithHalfOpenMaxCall returns a function to set the halfOpenMaxCall of Breaker
func WithHalfOpenMaxCall(d uint64) Option {
	return func(options *Breaker) {
		options.halfOpenMaxCall = d
	}
}

// WithStrategyOption returns a function to set the strategy function of a Breaker
func WithStrategyOption(o StrategyOption) Option {
	switch o.Strategy {
	case StrategyFail:
		if o.FailThreshold <= 0 {
			o.FailThreshold = DefaultFailThreshold
		}
		return func(options *Breaker) {
			options.strategyFn = FailStrategyFn(o.FailThreshold)
		}
	case StrategyContinuousFail:
		if o.ContinuousFailThreshold <= 0 {
			o.ContinuousFailThreshold = DefaultContinuousFailThreshold
		}
		return func(options *Breaker) {
			options.strategyFn = ContinuousFailStrategyFn(o.ContinuousFailThreshold)
		}
	case StrategyFailRate:
		if o.FailRate <= 0 || o.MinCall <= 0 {
			o.FailRate = DefaultFailRate
			o.MinCall = DefaultMinCall
		}
		return func(options *Breaker) {
			options.strategyFn = FailRateStrategyFn(o.FailRate, o.MinCall)
		}
	default:
		panic("unknown breaker strategy")
	}
}

// WithWebhook returns a function to set to notify of Breaker
func WithWebhook(webhook string) Option {
	return func(options *Breaker) {
		options.notify = notify.NewNotify(webhook)
	}
}
