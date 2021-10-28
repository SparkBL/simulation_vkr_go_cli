package main

import (
	"fmt"
	"rqsim/components"
	"sort"
)

func main() {
	conf, err := components.ParseConfig("conf.json")
	if err != nil {
		return
	}
	var inputChannel = make(chan components.Request, 1)
	var orbitChannel = make(chan components.Request, 1)
	var orbitAppendChannel = make(chan components.Request, 2)
	var outputChannel = make(chan components.Request, 1)
	var calledChannel = make(chan components.Request, 1)
	var inStream components.Process
	var sigmaDelay components.Delay
	switch conf.InputType {
	case "simple":
		inStream = components.NewSimpleStream(components.ExpDelay{Intensity: conf.LSimple}, components.TypeInput, inputChannel)
	case "mmpp":
		inStream = components.NewMMPP(conf.L, conf.Q, components.TypeInput, inputChannel)
	default:
		return
	}
	switch conf.SigmaDelayType {
	case "exp":
		sigmaDelay = components.ExpDelay{Intensity: conf.Sigma}
	case "uniform":
		sigmaDelay = components.UniformDelay{A: conf.SigmaA, B: conf.SigmaB}
	default:
		return
	}
	callStream := components.NewSimpleStream(components.ExpDelay{Intensity: conf.Alpha}, components.TypeCalled, calledChannel)
	orbit := components.NewOrbit(sigmaDelay, orbitChannel, orbitAppendChannel)
	node := components.NewNode(components.ExpDelay{Intensity: conf.Mu1}, components.ExpDelay{Intensity: conf.Mu2}, inputChannel, calledChannel, orbitChannel, orbitAppendChannel, outputChannel)
	components.Time = 0
	components.End = conf.End
	for components.Time < components.End {
		if len(components.EventQueue) > 0 {
			sort.Float64s(components.EventQueue)
			components.Time = components.EventQueue[0]
			components.EventQueue = components.EventQueue[1:]
		}
		node.Produce()
		orbit.Append()
		inStream.Produce()
		orbit.Produce()
		callStream.Produce()
		fmt.Printf("Time is %2f, Event Queue is %2f", components.Time, components.EventQueue)
		if len(outputChannel) > 0 {
			<-outputChannel
		}
	}
	fmt.Println(len(inputChannel), len(outputChannel), len(calledChannel), len(orbitAppendChannel), len(orbitChannel), len(outputChannel))
}
