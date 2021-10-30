package components

type Orbit struct {
	requests           []Request
	delay              Delay
	orbitChannel       chan Request
	orbitAppendChannel chan Request
	timeChangeChannel  chan bool
}

func (o *Orbit) Start() {
	go func() {
		for req := range o.orbitAppendChannel {
			req.StatusChangeAt = o.delay.Get()
			AppendToEventQueue(req.StatusChangeAt)
			req.Status = statusTravel
			o.requests = append(o.requests, req)
		}
	}()
	go func() {
		for range o.timeChangeChannel {
			if len(o.requests) > 0 {
				for _, v := range o.requests {
					if almostEqual(v.StatusChangeAt, Time) {
						ret := v
						o.requests = o.requests[1:]
						o.orbitChannel <- ret
						return
					}
				}

			}
		}
	}()
}

func NewOrbit(delay Delay, orbitChannel chan Request, orbitAppendChannel chan Request, TimeChangeChannel chan bool) *Orbit {
	return &Orbit{
		delay:              delay,
		orbitChannel:       orbitChannel,
		orbitAppendChannel: orbitAppendChannel,
		requests:           make([]Request, 0),
		timeChangeChannel:  TimeChangeChannel,
	}
}
