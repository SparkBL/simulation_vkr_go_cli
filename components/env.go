package components

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Process interface {
	Produce()
}

func Init() {
	Time = 0
}

var EventQueue []float64
var Time float64
var End float64
var Interval float64

type Config struct {
	InputType      string      `json:"input_type"`
	SigmaDelayType string      `json:"sigma_delay_type"`
	Sigma          float64     `json:"sigma"`
	SigmaA         float64     `json:"sigmaA"`
	SigmaB         float64     `json:"sigmaB"`
	L              [][]float64 `json:"L"`
	LSimple        float64     `json:"LSimple"`
	Q              [][]float64 `json:"Q"`
	Alpha          float64     `json:"alpha"`
	Mu1            float64     `json:"mu1"`
	Mu2            float64     `json:"mu2"`
	End            float64     `json:"end"`
	Interval       float64     `json:"interval"`
	Delimiter      rune        `json:"delimiter"`
}

func ParseConfig(configFile string) (Config, error) {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Couldn't open config file ", configFile)
		return Config{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf Config
	json.Unmarshal(byteValue, &conf)
	return conf, nil
}
