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
		FailThreshold: 3,
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

	breaker := NewBreaker(
		WithName("breakerName"),
		WithWindowInterval(10*time.Second),
		WithCoolDownTime(2*time.Second),
		WithHalfOpenMaxCall(2),
		WithStrategyOption(strategyFailOpt),
	)

	for i := 1; i <= 20; i++ {
		log.Println("i:", i)
		err := breaker.Do(func() error {
			rand.Seed(time.Now().UnixNano())
			if i == 1 || i == 2 || i == 3 || i == 6 || i == 11 {
				return errors.New("")
			} else {
				return nil
			}
		})
		if errors.Is(err, ErrStateHalfOpen) || errors.Is(err, ErrStateOpen) {
			log.Printf("err: %s", err.Error())
		}

		fmt.Println()
		time.Sleep(time.Second)
	}

}
