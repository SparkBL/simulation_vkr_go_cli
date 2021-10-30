package components

type IntervalStat struct {
	input  float64
	called float64
}

type StatCollector struct {
	intervalStats []IntervalStat
	outputChannel chan Request
	cur           IntervalStat
}

func (s *StatCollector) GatherStat() {
	for i := 0; i < len(s.outputChannel); i++ {
		if r := <-s.outputChannel; r.Type == TypeInput {
			s.cur.input++
		} else {
			s.cur.called++
		}
	}
}

func (s *StatCollector) ChangeInterval() {
	s.intervalStats = append(s.intervalStats, s.cur)
	s.cur = IntervalStat{input: 0, called: 0}
}

func NewStatCollector(outChannel chan Request) StatCollector {
	return StatCollector{
		intervalStats: make([]IntervalStat, 0),
		outputChannel: outChannel,
		cur:           IntervalStat{input: 0, called: 0},
	}
}

func (s *StatCollector) MeanInput() float64 {
	sum := 0.0
	for _, e := range s.intervalStats {
		sum += e.input
	}
	return sum / float64(len(s.intervalStats))
}

func (s *StatCollector) GetDistribution() [][]float64 {
	distSizeX, distSizeY := 0.0, 0.0
	for _, e := range s.intervalStats {
		if e.input > distSizeX {
			distSizeX = e.input
		}
		if e.called > distSizeY {
			distSizeY = e.called
		}
	}
	distr := make([][]float64, int(distSizeX+1))
	for i := range distr {
		distr[i] = make([]float64, int(distSizeY+1))
	}

	for _, e := range s.intervalStats {
		distr[int(e.input)][int(e.called)]++
	}
	for i := range distr {
		for j := range distr[i] {
			distr[i][j] /= float64(len(s.intervalStats))
		}
	}
	return distr
}
