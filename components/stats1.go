package components

/*type TimedStatCollector struct {
	distr         [][]float64
	outputChannel chan Request
	interval      float64
	curInterval   float64
}

func (s *TimedStatCollector) GatherStat() {
	for i := 0; i < len(s.outputChannel); i++ {
		r := <-s.outputChannel
		if !(r.StatusChangeAt>=Time-s.curInterval && r.StatusChangeAt<=s.curInterval){
			s.curInterval +=s.interval
		}
		for !(r.StatusChangeAt>=Time-s.curInterval && r.StatusChangeAt<=s.curInterval){
			s.distr[0][0] += s.interval
		}
		if r.Type == TypeInput {

		} else {
		}
	}
}

func NewTimedStatCollector(outChannel chan Request, interval float64) TimedStatCollector {
	return TimedStatCollector{
		distr:         make([][]float64, 10, 10),
		outputChannel: outChannel,
		interval:      interval,
		curInterval:   interval,
	}
}

func (s *StatCollector) TimedGetDistribution() [][]float64 {
	distSizeX, distSizeY := 0.0, 0.0
	for _, e := range s.intervalStats {
		if e.input > distSizeX {
			distSizeX = e.input
		}
		if e.called > distSizeY {
			distSizeY = e.called
		}
	}
	distr := make([][]float64, int(distSizeX), int(distSizeY))

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
*/
