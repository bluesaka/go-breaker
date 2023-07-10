package breaker

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

// TestBreaker
// exec `go test -v -run ^TestBreaker$` command to test
func TestBreaker(t *testing.T) {
	strategyFailOpt := StrategyOption{
		Strategy:      StrategyFail,
		FailThreshold: 2,
	}
	//strategyContinuousFailOpt := StrategyOption{
	//	Strategy:                StrategyContinuousFail,
	//	ContinuousFailThreshold: 2,
	//}
	//strategyFailRateOpt := StrategyOption{
	//	Strategy: StrategyFailRate,
	//	FailRate: 0.6,
	//	MinCall:  10,
	//}

	breaker := NewBreaker(WithName("breakerName"), WithWindowInterval(time.Second), WithCoolDownTime(time.Second),
		WithHalfOpenMaxCall(2), WithStrategyOption(strategyFailOpt))

	for i := 0; i < 20; i++ {
		log.Println("i:", i)
		breaker.Do(func() error {
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(2) == 1 {
				return nil
			} else {
				return errors.New("error")
			}
		})
		fmt.Println()
		time.Sleep(time.Millisecond * 500)
	}
}
