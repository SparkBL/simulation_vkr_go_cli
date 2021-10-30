package components

type MMPP struct {
	L                 [][]float64
	Q                 [][]float64
	RequestType       int
	state             int
	shiftTime         float64
	nextProduce       Request
	channel           chan Request
	timeChangeChannel chan float64
}

func (m *MMPP) shift(t float64) {
	if almostEqual(m.shiftTime, t) {
		sum, chance := 0.0, NextDouble()
		for i := 0; i < len(m.Q); i++ {
			if i != m.state {
				sum += m.Q[m.state][i] / (-m.Q[m.state][m.state])
				if chance <= sum {
					m.state = i
					m.nextProduce = Request{StatusChangeAt: ExponentialDelay(m.L[m.state][m.state]), Type: m.RequestType, Status: statusTravel}
					m.shiftTime = ExponentialDelay(-m.Q[m.state][m.state])
					AppendToEventQueue(m.nextProduce.StatusChangeAt)
					AppendToEventQueue(m.shiftTime)
				}
			}
		}
	}
}

func (m *MMPP) Start() {
	go func() {
		for t := range m.timeChangeChannel {
			m.shift(t)
			if almostEqual(m.nextProduce.StatusChangeAt, t) {
				m.channel <- m.nextProduce
				m.nextProduce = Request{StatusChangeAt: ExponentialDelay(m.L[m.state][m.state]), Type: m.RequestType, Status: statusTravel}
				AppendToEventQueue(m.nextProduce.StatusChangeAt)
			}
		}
	}()
}

type SimpleInput struct {
	nextProduce       Request
	delay             Delay
	RequestType       int
	channel           chan Request
	timeChangeChannel chan float64
}

func (s *SimpleInput) Start() {
	go func() {
		for t := range s.timeChangeChannel {
			if almostEqual(s.nextProduce.StatusChangeAt, t) {
				s.channel <- s.nextProduce
				s.nextProduce = Request{StatusChangeAt: s.delay.Get(), Type: s.RequestType, Status: statusTravel}
				AppendToEventQueue(s.nextProduce.StatusChangeAt)
			}
		}
	}()
}

func NewMMPP(L [][]float64, Q [][]float64, RequstType int, channel chan Request, TimeChangeChannel chan float64) *MMPP {
	nprod := Request{StatusChangeAt: ExponentialDelay(L[0][0]), Type: RequstType, Status: statusTravel}
	AppendToEventQueue(nprod.StatusChangeAt)
	//channel <- nprod
	return &MMPP{L: L,
		Q:                 Q,
		RequestType:       RequstType,
		state:             0,
		shiftTime:         ExponentialDelay(-Q[0][0]),
		nextProduce:       nprod,
		channel:           channel,
		timeChangeChannel: TimeChangeChannel,
	}
}

func NewSimpleStream(delay Delay, RequestType int, channel chan Request, TimeChangeChannel chan float64) *SimpleInput {
	nprod := Request{StatusChangeAt: delay.Get(), Type: RequestType, Status: statusTravel}
	//channel <- nprod
	AppendToEventQueue(nprod.StatusChangeAt)
	return &SimpleInput{
		nextProduce:       nprod,
		delay:             delay,
		RequestType:       RequestType,
		channel:           channel,
		timeChangeChannel: TimeChangeChannel,
	}
}
