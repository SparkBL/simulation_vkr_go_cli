package components

type Orbit struct {
	requests           []Request
	delay              Delay
	orbitChannel       *WriteRouter
	orbitAppendChannel *ReadRouter
}

func (o *Orbit) Append() {
	for i := 0; i < o.orbitAppendChannel.Len(); i++ {
		req := o.orbitAppendChannel.Pop()
		req.StatusChangeAt = o.delay.Get()
		EventQueue = append(EventQueue, req.StatusChangeAt)
		req.Status = statusTravel
		o.requests = append(o.requests, req)
	}

}

func (o *Orbit) Produce() {
	if len(o.requests) > 0 {
		for _, v := range o.requests {
			if almostEqual(v.StatusChangeAt, Time) {
				ret := v
				o.requests = o.requests[1:]
				o.orbitChannel.Push(ret)
				return
			}
		}

	}
}

func NewOrbit(delay Delay, orbitChannel *WriteRouter, orbitAppendChannel *ReadRouter) *Orbit {
	return &Orbit{
		delay:              delay,
		orbitChannel:       orbitChannel,
		orbitAppendChannel: orbitAppendChannel,
		requests:           make([]Request, 0),
	}
}
