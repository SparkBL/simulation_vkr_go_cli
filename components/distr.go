package components

import (
	"math"
	"math/rand"
)

const float64EqualityThreshold = 1e-3

func NextDouble() float64 {
	return rand.Float64()
}

type Delay interface {
	Get() float64
}

type ExpDelay struct {
	Intensity float64
}

func (e ExpDelay) Get() float64 {
	return rand.ExpFloat64()/e.Intensity + Time
}

type UniformDelay struct {
	A float64
	B float64
}

func (u UniformDelay) Get() float64 {
	return rand.Float64()*(u.B-u.A) + u.A + Time
}

func ExponentialDelay(intensity float64) float64 {
	return rand.ExpFloat64()/intensity + Time
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
