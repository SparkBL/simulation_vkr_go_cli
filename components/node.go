package components

type Node struct {
	nowServing         Request
	inputDelay         Delay
	calledDelay        Delay
	inChannel          *ReadRouter
	callChannel        *ReadRouter
	orbitChannel       *ReadRouter
	orbitAppendChannel *WriteRouter
	outChannel         chan<- Request
}

func (n *Node) Produce() {
	if n.nowServing.Status == statusServing && almostEqual(n.nowServing.StatusChangeAt, Time) {
		n.nowServing.Status = statusServed
		n.outChannel <- n.nowServing
	}
	if n.inChannel.Len() > 0 {
		if n.nowServing.Status == statusServing {
			n.orbitAppendChannel.Push(n.inChannel.Pop())
		} else {
			n.nowServing = n.inChannel.Pop()
			n.nowServing.StatusChangeAt, n.nowServing.Status = n.inputDelay.Get(), statusServing
			EventQueue = append(EventQueue, n.nowServing.StatusChangeAt)
		}
	}

	if n.orbitChannel.Len() > 0 {
		if n.nowServing.Status == statusServing {
			n.orbitAppendChannel.Push(n.orbitChannel.Pop())
		} else {
			n.nowServing = n.orbitChannel.Pop()
			n.nowServing.StatusChangeAt, n.nowServing.Status = n.inputDelay.Get(), statusServing
			EventQueue = append(EventQueue, n.nowServing.StatusChangeAt)
		}
	}

	if n.callChannel.Len() > 0 {
		if n.nowServing.Status != statusServing {
			n.nowServing = n.callChannel.Pop()
			n.nowServing.StatusChangeAt, n.nowServing.Status = n.calledDelay.Get(), statusServing
			EventQueue = append(EventQueue, n.nowServing.StatusChangeAt)
		} else {
			n.callChannel.Pop()
		}
	}
}
func NewNode(inputDelay Delay, calledDelay Delay, inChannel *ReadRouter, callChannel *ReadRouter, orbitChannel *ReadRouter, orbitAppendChannel *WriteRouter, outChannel chan<- Request) *Node {
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
