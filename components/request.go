package components

const (
	TypeInput  = 0
	TypeCalled = 1

	statusTravel  = 0
	statusServing = 1
	statusServed  = 2
)

type Request struct {
	Type           int
	Status         int
	StatusChangeAt float64
}
