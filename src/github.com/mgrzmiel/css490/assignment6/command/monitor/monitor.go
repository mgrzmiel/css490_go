// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel
// monitor periodically calls "/monitor" on the timeserver and authserver to collected statistics.
// It also print the statistics to the console as an json object

package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment6/lib/monitorConfig"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Sample struct {
	CurrentTime time.Time
	Value       int
}

type Statistics struct {
	Url        string
	Dictionery map[string][]Sample
}

// lock in ordert to print one object at the time
var lockPrinting *sync.RWMutex

// printResults
// This functon is responsible for printing the statistics to the console as an json object
func printResults(target string, monitorDictionary map[string][]Sample) {
	rates := marshalData(target, monitorDictionary)
	fmt.Printf("%s\n", rates)
}

// marshalData
// This function is responsible for marshelling the data
func marshalData(targetUrl string, monitorDictionary map[string][]Sample) string {
	statistics := Statistics{Url: targetUrl, Dictionery: monitorDictionary}
	data, err := json.Marshal(statistics)
	if err != nil {
		log.Errorf("Not able to marshall the data")
		return "Not able to marshall the data"
	} else {
		return string(data)
	}
}

// monitorTarget
// This function read the json object after from the webside endpoint /monitor every interval time
// When the timout passed, it prints the results.
func monitorTarget(target string) {
	timeout := time.Tick(time.Duration(monitorConfig.RuntimeSec) * time.Second)
	interval := time.Tick(time.Duration(monitorConfig.SampleIntervalSec) * time.Second)
	var monitorDictionary = make(map[string][]Sample)

	for {
		// if !strings.HasSuffix(target, "/") {
		// 	target += "/"
		// }

		requestURL := target + "monitor"
		client := http.Client{}
		res, err := client.Get(requestURL)
		if err != nil {
			log.Errorf("Get request error %s", err)
		} else {
			data, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Errorf("Error while reading the dirctionery: %s", err)
			} else {
				tempMap := make(map[string]int)
				err = json.Unmarshal([]byte(data), &tempMap)
				if err != nil {
					log.Warn("Unable to unmarshal the data")
				} else {
					timeNow := time.Now()
					for key, val := range tempMap {
						sample := Sample{CurrentTime: timeNow, Value: val}
						monitorDictionary[key] = append(monitorDictionary[key], sample)
					}
				}
			}
		}

		//wait the interval time to pass
		<-interval
		select {
		// if timeout passed, print the final statistics
		case <-timeout:
			lockPrinting.Lock()
			printResults(target, monitorDictionary)
			lockPrinting.Unlock()
			return
		// in other case continue
		default:
		}
	}
}

// main function
// This function is responsible for the flow of the program
func main() {
	logger, err := log.LoggerFromConfigAsFile(monitorConfig.LogPath)
	if err != nil {
		log.Errorf("Cannot open config file %s\n", err)
		return
	}
	log.ReplaceLogger(logger)

	// get the list of all targets
	targetsList := strings.Split(monitorConfig.Targets, ",")
	lockPrinting = new(sync.RWMutex)

	// for each target monitor the statistics
	for _, target := range targetsList {
		go monitorTarget(target)
	}

	// sleep until the monitoring will be done
	time.Sleep(time.Duration(2*monitorConfig.RuntimeSec) * time.Second)
}
