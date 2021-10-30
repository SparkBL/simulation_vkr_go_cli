package components

type Node struct {
	nowServing         Request
	inputDelay         Delay
	calledDelay        Delay
	inChannel          chan Request
	callChannel        chan Request
	orbitChannel       chan Request
	orbitAppendChannel chan Request
	outChannel         chan Request
	timeChangeChannel  chan float64
}

func (n *Node) Start() {
	go func() {
		for t := range n.timeChangeChannel {
			if n.nowServing.Status == statusServing && almostEqual(n.nowServing.StatusChangeAt, t) {
				n.nowServing.Status = statusServed
				n.outChannel <- n.nowServing
			}
		}
	}()
	go func() {
		for r := range n.inChannel {
			if n.nowServing.Status == statusServing {
				n.orbitAppendChannel <- r
			} else {
				n.nowServing = r
				n.nowServing.StatusChangeAt, n.nowServing.Status = n.inputDelay.Get(), statusServing
				AppendToEventQueue(n.nowServing.StatusChangeAt)
			}
		}
	}()

	go func() {
		for r := range n.orbitChannel {
			if n.nowServing.Status == statusServing {
				n.orbitAppendChannel <- r
			} else {
				n.nowServing = r
				n.nowServing.StatusChangeAt, n.nowServing.Status = n.inputDelay.Get(), statusServing
				AppendToEventQueue(n.nowServing.StatusChangeAt)
			}
		}
	}()

	go func() {
		for r := range n.callChannel {
			if n.nowServing.Status != statusServing {
				n.nowServing = r
				n.nowServing.StatusChangeAt, n.nowServing.Status = n.calledDelay.Get(), statusServing
				AppendToEventQueue(n.nowServing.StatusChangeAt)
			}
		}
	}()
}
func NewNode(inputDelay Delay, calledDelay Delay, inChannel chan Request, callChannel chan Request, orbitChannel chan Request, orbitAppendChannel chan Request, outChannel chan Request, TimeChangeChannel chan float64) *Node {
	return &Node{inputDelay: inputDelay,
		calledDelay:        calledDelay,
		inChannel:          inChannel,
		callChannel:        callChannel,
		orbitChannel:       orbitChannel,
		orbitAppendChannel: orbitAppendChannel,
		outChannel:         outChannel,
		nowServing:         Request{Status: statusServed},
		timeChangeChannel:  TimeChangeChannel,
	}
}
