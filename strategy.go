package breaker

const (
	StrategyFail           = iota + 1 // 失败数策略
	StrategyContinuousFail            // 连续失败数策略
	StrategyFailRate                  // 失败率策略
)

// StrategyFn breaker strategy function on request error
type StrategyFn func(metric Metric) bool

// FailStrategyFn breaker strategy function based on fail count
func FailStrategyFn(threshold uint64) StrategyFn {
	return func(metric Metric) bool {
		return metric.TotalFail >= threshold
	}
}

// ContinuousFailStrategyFn breaker strategy function based on continuous fail count
func ContinuousFailStrategyFn(threshold uint64) StrategyFn {
	return func(metric Metric) bool {
		return metric.ContinuousFail >= threshold
	}
}

// ContinuousFailStrategyFn breaker strategy function based on fail rate
func FailRateStrategyFn(rate float64, minCall uint64) StrategyFn {
	return func(metric Metric) bool {
		if metric.TotalRequest < minCall {
			return false
		}
		return float64(metric.TotalFail)/float64(metric.TotalRequest) >= rate
	}
}
