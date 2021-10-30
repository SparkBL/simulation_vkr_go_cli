package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"rqsim/components"
	"sort"
	"strconv"
	"time"
)

func writeToCSV(outputFile string, rows [][]float64) {
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	rowsStr := make([][]string, len(rows))
	for i := range rowsStr {
		rowsStr[i] = make([]string, len(rows[0]))
	}
	for i := range rows {
		for j := range rows[i] {
			rowsStr[i][j] = strconv.FormatFloat(rows[i][j], 'f', 6, 64)
		}
	}
	writer := csv.NewWriter(f)
	writer.Comma = ';'

	writer.WriteAll(rowsStr)

	f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	configFile := flag.String("c", "conf.json", "Path to json config file of rq system")
	outputFile := flag.String("o", "out.csv", "name of output file")
	flag.Parse()
	conf, err := components.ParseConfig(*configFile)
	if err != nil {
		return
	}
	var inputChannel = make(chan components.Request, 1)
	var orbitChannel = make(chan components.Request, 1)
	var orbitAppendChannel = make(chan components.Request, 10)
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
	statCollector := components.NewStatCollector(outputChannel)
	components.Time = 0
	components.End = conf.End
	components.Interval = conf.Interval
	go func() {
		for {
			fmt.Printf("Simulating for %2f. End at %2f\r", components.Time, components.End)
			time.Sleep(time.Second)
		}
	}()
	go statCollector.GatherStat()
	for components.Time < components.End {
		inStream.Produce()
		orbit.Append()
		orbit.Produce()
		callStream.Produce()
		node.Produce()
		if len(components.EventQueue) > 0 {
			sort.Float64s(components.EventQueue)
			components.Time, components.EventQueue = components.EventQueue[0], components.EventQueue[1:]
		}
	}
	close(outputChannel)

	writeToCSV(*outputFile, statCollector.GetDistribution())
}
