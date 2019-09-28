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
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"text/tabwriter"
	"time"
)

type step struct {
	setpoint float64
	input    float64
	output   float64
}

var tests = []struct {
	name         string
	p            float64
	i            float64
	d            float64
	min          float64
	max          float64
	steps        []*step
	stepDuration time.Duration
	expectPanic  bool
}{
	{
		name: "empty controller",
		steps: []*step{
			{0, 0, 0},
		},
	},
	{
		name: "p-only controller",
		p:    0.5,
		steps: []*step{
			{10, 5, 2.5},
			{0, 10, 0},
			{0, 15, -2.5},
			{0, 100, -45},
			{1, 0, 0.5},
		},
	},
	{
		name:         "i-only controller",
		i:            0.5,
		stepDuration: time.Second,
		steps: []*step{
			{10, 5, 2.5},
			{0, 5, 5},
			{0, 5, 7.5},
			{0, 15, 5},
			{0, 20, 0},
		},
	},
	{
		name:         "d-only controller",
		d:            0.5,
		stepDuration: time.Second,
		steps: []*step{
			{10, 5, -2.5},
			{0, 5, 0},
			{0, 10, -2.5},
		},
	},
	{
		name:         "pid controller",
		p:            0.5,
		i:            0.5,
		d:            0.5,
		stepDuration: time.Second,
		steps: []*step{
			{10, 5, 2.5},
			{0, 10, 0},
			{0, 15, -5},
			{0, 100, -132.5},
			{1, 0, 6},
		},
	},
	{
		name:         "thermostat example",
		p:            0.6,
		i:            1.2,
		d:            0.075,
		max:          1,
		stepDuration: time.Second,
		steps: []*step{
			{72, 50, 1},
			{0, 51, 1},
			{0, 55, 1},
			{0, 60, 1},
			{0, 75, 0},
			{0, 76, 0},
			{0, 74, 0},
			{0, 72, 0.15},
			{0, 71, 1},
		},
	},
	{
		name: "pd controller",
		p:    40.0,
		i:    0,
		d:    12.0,
		max:  255, min: 0,
		stepDuration: time.Second,
		steps: []*step{
			{90.0, 22.00, 255.00},
			{0, 25.29, 255.00},
			{0, 28.56, 255.00},
			{0, 31.80, 255.00},
			{0, 35.02, 255.00},
			{0, 38.21, 255.00},
			{0, 41.38, 255.00},
			{0, 44.53, 255.00},
			{0, 47.66, 255.00},
			{0, 50.76, 255.00},
			{0, 53.84, 255.00},
			{0, 56.90, 255.00},
			{0, 59.93, 255.00},
			{0, 62.95, 255.00},
			{0, 65.94, 255.00},
			{0, 68.91, 255.00},
			{0, 71.85, 255.00},
			{0, 74.78, 255.00},
			{0, 77.69, 255.00},
			{0, 80.57, 255.00},
			{0, 83.43, 228.48},
			{0, 85.93, 132.80},
			{0, 87.18, 97.80},
			{0, 87.96, 72.24},
			{0, 88.41, 58.20},
			{0, 88.68, 49.56},
			{0, 88.83, 45.00},
			{0, 88.92, 42.12},
			{0, 88.98, 40.08},
			{0, 89.00, 39.76},
			{0, 89.03, 38.44},
			{0, 89.03, 38.80},
			{0, 89.05, 37.76},
			{0, 89.04, 38.52},
			{0, 89.05, 37.88},
			{0, 89.05, 38.00},
			{0, 89.05, 38.00},
			{0, 89.05, 38.00},
			{0, 89.05, 38.00},
		},
	},
}

func TestController(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case r == nil:
				case test.expectPanic && fmt.Sprintf("%T", r) == "*pidctl.MinMaxError":
					t.Log("trapped:", r)
				default:
					panic(r)
				}
			}()
			c := Controller{}
			if test.p != 0 {
				c.P = ratFloat64(test.p)
			}
			if test.i != 0 {
				c.I = ratFloat64(test.i)
			}
			if test.d != 0 {
				c.D = ratFloat64(test.d)
			}
			if test.min != 0 || test.max != 0 {
				c.Min = ratFloat64(test.min)
				c.Max = ratFloat64(test.max)
			}

			var buf bytes.Buffer
			log := tabwriter.NewWriter(&buf, 8, 0, 1, ' ', 0)
			fmt.Fprint(log, "\tcycle\tgot\texpected\tsetpoint\tinput\toutput\n")
			for i, u := range test.steps {
				if u.setpoint != 0 {
					c.Setpoint = ratFloat64(u.setpoint)
				}
				got := c.Accumulate(ratFloat64(u.input), test.stepDuration)
				roundedGot, roundedExpected := got.FloatString(3), ratFloat64(u.output).FloatString(3)
				msg := ""
				if roundedGot != roundedExpected {
					msg = "error"
					t.Fail()
				}
				fmt.Fprintf(log, "%s\t%d\t%v\t%v\t%v\t%v\t%v\n", msg, i, roundedGot, roundedExpected, u.setpoint, u.input, u.output)
			}
			log.Flush()
			scanner := bufio.NewScanner(&buf)
			for scanner.Scan() {
				t.Log(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				panic(fmt.Sprint("reading table output:", err))
			}
		})
	}
}

func ratFloat64(i float64) *big.Rat {
	return new(big.Rat).SetFloat64(i)
}
