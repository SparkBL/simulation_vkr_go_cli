package components

type MMPP struct {
	L           [][]float64
	Q           [][]float64
	RequestType int
	state       int
	shiftTime   float64
	nextProduce Request
	channel     chan<- Request
}

func (m *MMPP) shift() {
	if almostEqual(m.shiftTime, Time) {
		sum, chance := 0.0, rng.Float64()
		for i := 0; i < len(m.Q); i++ {
			if i != m.state {
				sum += m.Q[m.state][i] / (-m.Q[m.state][m.state])
				if chance <= sum {
					m.state = i
					m.nextProduce = Request{StatusChangeAt: ExponentialDelay(m.L[m.state][m.state]), Type: m.RequestType, Status: statusTravel}
					m.shiftTime = ExponentialDelay(-m.Q[m.state][m.state])
					EventQueue = append(EventQueue, m.nextProduce.StatusChangeAt)
					EventQueue = append(EventQueue, m.shiftTime)
				}
			}
		}
	}
}

func (m *MMPP) Produce() {
	m.shift()
	if almostEqual(m.nextProduce.StatusChangeAt, Time) {
		m.channel <- m.nextProduce
		m.nextProduce = Request{StatusChangeAt: ExponentialDelay(m.L[m.state][m.state]), Type: m.RequestType, Status: statusTravel}
		EventQueue = append(EventQueue, m.nextProduce.StatusChangeAt)
	}
}

type SimpleInput struct {
	nextProduce Request
	delay       Delay
	RequestType int
	channel     chan<- Request
}

func (s *SimpleInput) Produce() {
	if almostEqual(s.nextProduce.StatusChangeAt, Time) {
		s.channel <- s.nextProduce
		s.nextProduce = Request{StatusChangeAt: s.delay.Get(), Type: s.RequestType, Status: statusTravel}
		EventQueue = append(EventQueue, s.nextProduce.StatusChangeAt)
	}
}

func NewMMPP(L [][]float64, Q [][]float64, RequstType int, channel chan<- Request) *MMPP {
	nprod := Request{StatusChangeAt: ExponentialDelay(L[0][0]), Type: RequstType, Status: statusTravel}
	EventQueue = append(EventQueue, nprod.StatusChangeAt)
	//channel <- nprod
	return &MMPP{L: L,
		Q:           Q,
		RequestType: RequstType,
		state:       0,
		shiftTime:   ExponentialDelay(-Q[0][0]),
		nextProduce: nprod,
		channel:     channel}
}

func NewSimpleStream(delay Delay, RequestType int, channel chan<- Request) *SimpleInput {
	nprod := Request{StatusChangeAt: delay.Get(), Type: RequestType, Status: statusTravel}
	//channel <- nprod
	EventQueue = append(EventQueue, nprod.StatusChangeAt)
	return &SimpleInput{
		nextProduce: nprod,
		delay:       delay,
		RequestType: RequestType,
		channel:     channel}
}
