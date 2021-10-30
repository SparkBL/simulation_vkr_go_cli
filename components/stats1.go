package components

type TimedStatCollector struct {
	distr         [][]float64
	outputChannel chan Request
	curCount      IntervalStat
}

func (s *TimedStatCollector) GatherStat() {
	curInterval := Time + Interval
	curPoint := 0.0
	for r := range s.outputChannel {
		for r.StatusChangeAt > curInterval {
			s.distr[s.curCount.input][s.curCount.called] += curInterval - curPoint
			curPoint = curInterval
			curInterval += Interval
			s.curCount = IntervalStat{input: 0, called: 0}
		}
		s.distr[s.curCount.input][s.curCount.called] += r.StatusChangeAt - curPoint
		curPoint = r.StatusChangeAt
		switch r.Type {
		case TypeInput:
			s.curCount.input++
		case TypeCalled:
			s.curCount.called++
		}
	}
}

func NewTimedStatCollector(outChannel chan Request) TimedStatCollector {
	distrw := make([][]float64, 20)
	for i := range distrw {
		distrw[i] = make([]float64, 2)
	}
	return TimedStatCollector{
		distr:         distrw,
		outputChannel: outChannel,
		curCount:      IntervalStat{input: 0, called: 0},
	}
}

func (s *TimedStatCollector) GetDistribution() [][]float64 {
	normalize := 0.0
	for i := 0; i < len(s.distr); i++ {
		for j := 0; j < len(s.distr[i]); j++ {
			normalize += s.distr[i][j]
		}
	}
	for i := range s.distr {
		for j := range s.distr[i] {
			s.distr[i][j] /= normalize
		}
	}
	return s.distr
}
