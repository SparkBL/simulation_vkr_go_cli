package components

type TimedStatCollector struct {
	distr         [][]float64
	outputChannel <-chan Request
	maxInput      int
	maxCalled     int
}

func (s *TimedStatCollector) GatherStat() {
	curInputCount, curCalledCount := 0, 0
	prevInput, prevCalled := 0.0, 0.0
	curCalledQueue, curInputQueue := []float64{}, []float64{}
	for r := range s.outputChannel {
		switch r.Type {
		case TypeInput:
			for len(curInputQueue) > 0 && curInputQueue[0] < r.StatusChangeAt {
				s.distr[curInputCount][curCalledCount] += curInputQueue[0] - prevInput
				curInputCount--
				prevInput, curInputQueue = curInputQueue[0], curInputQueue[1:]
			}
			s.distr[curInputCount][curCalledCount] += r.StatusChangeAt - prevInput
			curInputCount++
			curInputQueue, prevInput = append(curInputQueue, r.StatusChangeAt+Interval), r.StatusChangeAt
			if curInputCount > s.maxInput {
				s.maxInput = curInputCount
			}
		case TypeCalled:
			for len(curCalledQueue) > 0 && curCalledQueue[0] < r.StatusChangeAt {
				s.distr[curInputCount][curCalledCount] += curCalledQueue[0] - prevCalled
				curCalledCount--
				prevCalled, curCalledQueue = curCalledQueue[0], curCalledQueue[1:]
			}
			s.distr[curInputCount][curCalledCount] += r.StatusChangeAt - prevCalled
			curCalledCount++
			curCalledQueue, prevCalled = append(curCalledQueue, r.StatusChangeAt+Interval), r.StatusChangeAt
			if curCalledCount > s.maxCalled {
				s.maxCalled = curCalledCount
			}
		}
	}
}

func NewTimedStatCollector(outChannel <-chan Request) TimedStatCollector {
	distrw := make([][]float64, 30)
	for i := range distrw {
		distrw[i] = make([]float64, 30)
	}
	return TimedStatCollector{
		distr:         distrw,
		outputChannel: outChannel,
	}
}

func (s *TimedStatCollector) GetDistribution() [][]float64 {
	normalize := 0.0
	for i := 0; i < len(s.distr); i++ {
		for j := 0; j < len(s.distr[0]); j++ {
			normalize += s.distr[i][j]
		}
	}
	ret := make([][]float64, s.maxInput)
	for i := range ret {
		ret[i] = make([]float64, s.maxCalled)
	}
	for i := 0; i < s.maxInput; i++ {
		for j := 0; j < s.maxCalled; j++ {
			ret[i][j] = s.distr[i][j] / normalize
		}
	}
	return ret
}
