package components

type Node struct {
	nowServing         Request
	inputDelay         Delay
	calledDelay        Delay
	inChannel          <-chan Request
	callChannel        <-chan Request
	orbitChannel       <-chan Request
	orbitAppendChannel chan<- Request
	outChannel         chan<- Request
}

func (n *Node) Produce() {
	if n.nowServing.Status == statusServing && almostEqual(n.nowServing.StatusChangeAt, Time) {
		n.nowServing.Status = statusServed
		n.outChannel <- n.nowServing
	}
	if len(n.inChannel) > 0 {
		if n.nowServing.Status == statusServing {
			n.orbitAppendChannel <- <-n.inChannel
		} else {
			n.nowServing = <-n.inChannel
			n.nowServing.StatusChangeAt, n.nowServing.Status = n.inputDelay.Get(), statusServing
			EventQueue = append(EventQueue, n.nowServing.StatusChangeAt)
		}
	}

	if len(n.orbitChannel) > 0 {
		if n.nowServing.Status == statusServing {
			n.orbitAppendChannel <- <-n.orbitChannel
		} else {
			n.nowServing = <-n.orbitChannel
			n.nowServing.StatusChangeAt, n.nowServing.Status = n.inputDelay.Get(), statusServing
			EventQueue = append(EventQueue, n.nowServing.StatusChangeAt)
		}
	}

	if len(n.callChannel) > 0 {
		if n.nowServing.Status != statusServing {
			n.nowServing = <-n.callChannel
			n.nowServing.StatusChangeAt, n.nowServing.Status = n.calledDelay.Get(), statusServing
			EventQueue = append(EventQueue, n.nowServing.StatusChangeAt)
		} else {
			<-n.callChannel
		}
	}
}
func NewNode(inputDelay Delay, calledDelay Delay, inChannel <-chan Request, callChannel <-chan Request, orbitChannel <-chan Request, orbitAppendChannel chan<- Request, outChannel chan<- Request) *Node {
	return &Node{inputDelay: inputDelay,
		calledDelay:        calledDelay,
		inChannel:          inChannel,
		callChannel:        callChannel,
		orbitChannel:       orbitChannel,
		orbitAppendChannel: orbitAppendChannel,
		outChannel:         outChannel,
		nowServing:         Request{Status: statusServed},
	}
}
