package breaker

import "time"

// Metric 指标
type Metric struct {
	WindowBatch       uint64    // 窗口批号
	WindowExpiry      time.Time // 窗口结束时间
	TotalRequest      uint64    // 总请求数
	TotalSuccess      uint64    // 总成功数
	TotalFail         uint64    // 总失败数
	ContinuousSuccess uint64    // 连续成功数
	ContinuousFail    uint64    // 连续失败数
}

// NewWindowBatch new window batch
func (m *Metric) NewWindowBatch() {
	m.WindowBatch++
}

// onRequest on request
func (m *Metric) onRequest() {
	m.TotalRequest++
}

// onSuccess on success call
func (m *Metric) onSuccess() {
	m.TotalSuccess++
	m.ContinuousSuccess++
	m.ContinuousFail = 0
}

// onFail on fail call
func (m *Metric) onFail() {
	m.TotalFail++
	m.ContinuousFail++
	m.ContinuousSuccess = 0
}

// OnReset reset window
func (m *Metric) OnReset() {
	m.TotalRequest = 0
	m.TotalSuccess = 0
	m.TotalFail = 0
	m.ContinuousSuccess = 0
	m.ContinuousFail = 0
}
