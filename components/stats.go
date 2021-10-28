package components

type IntervalStat struct {
	input  float64
	called float64
}

type StatCollector struct {
	intervalStats        []IntervalStat
	outputChannel        chan Request
	changeIntervalCannel chan bool
	cur                  IntervalStat
}

func (s *StatCollector) GatherStat() {
	cur := IntervalStat{input: 0, called: 0}
	for r := range s.outputChannel {
		if r.Type == TypeInput {
			cur.input++
		} else {
			cur.called++
		}
	}
}

func (s *StatCollector) ChangeInterval() {
	s.intervalStats = append(s.intervalStats)
}

func NewStatCollector(outChannel chan Request, changeIntervalCannel chan bool) StatCollector {
	return StatCollector{
		intervalStats:        make([]IntervalStat, 0),
		outputChannel:        outChannel,
		changeIntervalCannel: changeIntervalCannel,
	}
}
