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

package pidctl

import (
	"math/big"
	"sync"
	"time"
)

var ratZero = new(big.Rat).SetInt64(0)

// Controller implements a PID controller.
type Controller struct {
	// P is the proportional gain
	P *big.Rat
	// I is integral gain
	I *big.Rat
	// D is the derivative gain
	D *big.Rat
	// current setpoint
	Setpoint *big.Rat
	// Min the lowest value acceptable for the Output
	Min *big.Rat
	// Max the lowest value acceptable for the Output
	Max *big.Rat

	prevValue *big.Rat
	integral  *big.Rat
	initOnce  sync.Once
}

func (p *Controller) init() {
	p.initOnce.Do(func() {
		if p.P == nil {
			p.P = new(big.Rat).SetInt64(0)
		}
		if p.I == nil {
			p.I = new(big.Rat).SetInt64(0)
		}
		if p.D == nil {
			p.D = new(big.Rat).SetInt64(0)
		}
		if p.Setpoint == nil {
			p.Setpoint = new(big.Rat).SetInt64(0)
		}
		if p.prevValue == nil {
			p.prevValue = new(big.Rat).SetInt64(0)
		}
		if p.integral == nil {
			p.integral = new(big.Rat).SetInt64(0)
		}
	})
}

// Accumulate updates the controller with the given value and duration since the
// last update. It returns the new output.
func (p *Controller) Accumulate(v *big.Rat, duration time.Duration) *big.Rat {
	p.init()
	var (
		value      = v.Set(v)
		dt         = new(big.Rat)
		err        = new(big.Rat)
		integral   = new(big.Rat)
		derivative = new(big.Rat)
		output     = new(big.Rat)
		prevValue  *big.Rat
	)
	prevValue, p.prevValue = p.prevValue, value

	dt.SetInt64(int64(duration / time.Second))
	err.Sub(p.Setpoint, value)
	integral.
		Mul(err, dt).
		Mul(integral, p.I).
		Add(integral, p.integral)
	p.integral = integral
	if p.Max != nil && p.integral.Cmp(p.Max) > 0 {
		p.integral = p.Max
	} else if p.Min != nil && p.integral.Cmp(p.Min) < 0 {
		p.integral = p.Min
	}

	if dt.Cmp(ratZero) > 0 {
		derivative.Sub(value, prevValue)
		derivative.Quo(derivative, dt)
		derivative.Neg(derivative)
	}

	output.
		Add(output, err.Mul(p.P, err)).
		Add(output, p.integral).
		Add(output, derivative.Mul(derivative, p.D))

	if p.Max != nil && output.Cmp(p.Max) > 0 {
		output = p.Max
	} else if p.Min != nil && output.Cmp(p.Min) < 0 {
		output = p.Min
	}

	return output
}
