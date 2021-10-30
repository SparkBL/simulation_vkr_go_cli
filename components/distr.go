package components

import (
	"math"
	"math/rand"
	"time"

	"github.com/seehuhn/mt19937"
)

const float64EqualityThreshold = 1e-14

var rng = rand.New(mt19937.New())

func GenSeed() {
	rng.Seed(time.Now().UnixNano())
}

func NextDouble() float64 {
	return rng.Float64()
}

type Delay interface {
	Get() float64
}

type ExpDelay struct {
	Intensity float64
}

func (e ExpDelay) Get() float64 {
	return rng.ExpFloat64()/e.Intensity + Time
}

type UniformDelay struct {
	A float64
	B float64
}

func (u UniformDelay) Get() float64 {
	return rng.Float64()*(u.B-u.A) + u.A + Time
}

func ExponentialDelay(intensity float64) float64 {
	return rng.ExpFloat64()/intensity + Time
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
