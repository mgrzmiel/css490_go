// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Sample struct {
	currentTime time.Time
	value       int
}

var targets string
var sampleIntervalSec int
var runtimeSec int
var lockPrinting *sync.RWMutex

func printResults(target string, monitorDictionary map[string][]Sample) {
	fmt.Printf("url: \t%s\n", target)
	for key, value := range monitorDictionary {
		fmt.Printf("key: \t%s \t\n", key)
		for _, oneVal := range value {
			fmt.Printf("time: \t%s\t", oneVal.currentTime.Format("3:04:05 PM"))
			fmt.Printf("value: \t%d\t\n", oneVal.value)
		}
	}
}

func getRates(target string, monitorDictionary map[string][]Sample) {
	ratesMap := make(map[string]float64)
	for key, value := range monitorDictionary {
		length := len(value)
		if length < 2 {
			return
		} else {
			lastIndex := length - 1
			secondLastIndex := length - 2
			time1 := value[lastIndex].currentTime
			time2 := value[secondLastIndex].currentTime
			// fmt.Printf("time1: \t%s \t\n", time1)
			// fmt.Printf("key: \t%s \t\n", key)
			timeDiff := time1.Sub(time2)
			//timeDiff := value[lastIndex].currentTime - value[secondLastIndex].currentTime
			countDiff := value[lastIndex].value - value[secondLastIndex].value
			rate := float64(countDiff) / timeDiff.Seconds()
			ratesMap[key] = rate
		}
	}
	fmt.Printf("url: \t%s\n", target)
	for key, value := range ratesMap {
		fmt.Printf("key: \t%s \t", key)
		fmt.Printf("rate: \t%f\t\n", value)
	}
}

func monitorTarget(target string) {

	timeout := time.Tick(time.Duration(runtimeSec) * time.Second)
	interval := time.Tick(time.Duration(sampleIntervalSec) * time.Second)
	var monitorDictionary = make(map[string][]Sample)

	for {
		requestURL := target + "/monitor"
		client := http.Client{}
		res, err := client.Get(requestURL)
		if err != nil {
			//log.Errorf("Get request error %s", err)
		} else {
			data, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				//log.Errorf("Error while reading the dirctionery: %s", err)
			} else {
				tempMap := make(map[string]int)
				err = json.Unmarshal([]byte(data), &tempMap)
				if err != nil {
					//log.Warn("Unable to unmarshal the data")
				} else {
					timeNow := time.Now()
					for key, val := range tempMap {
						sample := Sample{currentTime: timeNow, value: val}
						monitorDictionary[key] = append(monitorDictionary[key], sample)
					}
				}
			}
		}
		lockPrinting.Lock()
		//printResults(target, monitorDictionary)
		getRates(target, monitorDictionary)
		lockPrinting.Unlock()
		<-interval

		select {
		case <-timeout:
			return
		default:
		}
	}

}

func main() {
	flag.StringVar(&targets, "targets", "", "rate")
	flag.IntVar(&sampleIntervalSec, "sample-interval-sec", 0, "sample-interval-sec")
	flag.IntVar(&runtimeSec, "runtime-sec", 0, "runtime-sec")

	flag.Parse()

	targetsList := strings.Split(targets, ",")
	lockPrinting = new(sync.RWMutex)

	for _, target := range targetsList {
		go monitorTarget(target)
	}

	time.Sleep(time.Duration(2*runtimeSec) * time.Second)
}
