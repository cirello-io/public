/*
Copyright 2019 github.com/ucirello and cirello.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pidctl_test

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"cirello.io/pidctl"
)

// A car whose PID controller is going to act to try to stabilize its
// speed to the given setpoint.
func Example() {
	car := pidctl.Controller{
		P:        big.NewRat(1, 5),
		I:        big.NewRat(1, 100),
		D:        big.NewRat(1, 15),
		Min:      big.NewRat(-1, 2), // min acceleration rate: 0.5 mps
		Max:      big.NewRat(10, 2), // max acceleration rate: 5 mps
		Setpoint: big.NewRat(60, 1), // target speed: 60 mph
	}
	speed := float64(20) // the car starts in motion. 20mph
	const travel = 15 * time.Second
	for i := time.Second; i <= travel; i += time.Second {
		desiredThrottle := car.Accumulate(new(big.Rat).SetFloat64(speed), time.Second)
		actualThrottle, _ := desiredThrottle.Float64()
		actualThrottle = math.Ceil(actualThrottle)
		fmt.Printf("%s speed: %.2f throttle: %.2f (desired: %s)\n", i, speed, actualThrottle, desiredThrottle.FloatString(2))
		speed += actualThrottle
		switch i % 5 {
		case 0:
			// head wind starts strong: 2mps
			speed -= 2
		case 1:
			// head wind ends weak: 1mps
			speed -= 1
		}
	}

	// Output:
	// 1s speed: 20.00 throttle: 5.00 (desired: 5.00)
	// 2s speed: 23.00 throttle: 5.00 (desired: 5.00)
	// 3s speed: 26.00 throttle: 5.00 (desired: 5.00)
	// 4s speed: 29.00 throttle: 5.00 (desired: 5.00)
	// 5s speed: 32.00 throttle: 5.00 (desired: 5.00)
	// 6s speed: 35.00 throttle: 5.00 (desired: 5.00)
	// 7s speed: 38.00 throttle: 5.00 (desired: 5.00)
	// 8s speed: 41.00 throttle: 5.00 (desired: 5.00)
	// 9s speed: 44.00 throttle: 5.00 (desired: 5.00)
	// 10s speed: 47.00 throttle: 5.00 (desired: 5.00)
	// 11s speed: 50.00 throttle: 5.00 (desired: 4.55)
	// 12s speed: 53.00 throttle: 5.00 (desired: 4.02)
	// 13s speed: 56.00 throttle: 4.00 (desired: 3.46)
	// 14s speed: 58.00 throttle: 4.00 (desired: 3.15)
	// 15s speed: 60.00 throttle: 3.00 (desired: 2.75)
}
