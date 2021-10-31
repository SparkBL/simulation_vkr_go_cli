package components

type ReadRouter struct {
	q *[]Request
}

type WriteRouter struct {
	q *[]Request
}

func (r *ReadRouter) Pop() Request {
	ret := (*r.q)[0]
	*r.q = (*r.q)[1:]
	return ret
}

func (r *ReadRouter) Len() int {
	return len(*r.q)
}

func (w *WriteRouter) Push(e Request) {
	*w.q = append(*w.q, e)
}

func NewRouter() (*ReadRouter, *WriteRouter) {
	storage := make([]Request, 0)
	return &ReadRouter{q: &storage}, &WriteRouter{q: &storage}
}
