// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel
// SessionManager is resposnisible for managing the session. 

package sessionManager

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// Session structure keeps the map with the session keys and
// RWMutex object for the synchronization purpose
type Sessions struct {
	SessionsMap     map[string]string
	SessionsSyncLoc *sync.RWMutex
}

// New
// Creates new Session structure
// If there is a file with data, read the data from file to map
// Otherwise reutrn empty map
func New(fileName string) *Sessions {
	sessions := Sessions{SessionsMap: make(map[string]string), SessionsSyncLoc: new(sync.RWMutex)}
	if _, err := os.Stat(fileName); err == nil {
		sessions.SessionsMap = readFromFile(fileName)
	}
	return &sessions
}

// readFromFile
// It returns the map with data if the file exist and the data can be unmarshal
func readFromFile(fileName string) map[string]string {
	tempMap := make(map[string]string)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Infof("Cannot read from file: %s\n", fileName)
		return nil
	} else {
		err = json.Unmarshal([]byte(file), &tempMap)
		if err != nil {
			log.Warnf("Unable to unmarshal the data from file: %s, \n", fileName)
			return nil
		} else {
			return tempMap
		}
	}
}

// WriteToFile
// This function is performed every given interval time.
// It rename the existing file, copying the data from map to another map
// and then save the data from that map to file.
// After that it make sure the data was saved correctly
func (s *Sessions) WriteToFile(fileName string, checkpointInterval time.Duration) {
	performWriting := time.Tick(checkpointInterval * time.Second)
	for _ = range performWriting {
		// rename the file if it exists
		newFileName := fileName + ".bak"
		if _, err := os.Stat(fileName); err == nil {
			err := os.Rename(fileName, newFileName)
			if err != nil {
				log.Errorf("Not able to rename the file")
			}
		}

		// write to temp dictionery
		tempMap := make(map[string]string)
		s.SessionsSyncLoc.Lock()
		for key, value := range s.SessionsMap {
			tempMap[key] = value
		}
		s.SessionsSyncLoc.Unlock()

		// marshall the data
		data, err := json.Marshal(tempMap)
		if err != nil {
			log.Warn("Not able to marshall the data")
		}
		// write to file
		err = ioutil.WriteFile(fileName, data, 0644)
		if err != nil {
			log.Warn("Not able to save the file")
		}

		// compare the content of the saved file with the dictionery
		fileContent := readFromFile(fileName)
		equal := compareMaps(tempMap, fileContent)

		// if the saved copy of data is correct, removed old file 
		if equal {
			err = os.Remove(newFileName)
			if err != nil {
				log.Errorf("Not able to remove the backfile: %s", newFileName)
			}
		// otherwise restore the old file
		} else {
			err := os.Rename(newFileName, fileName)
			if err != nil {
				log.Error("Not able to rename the file")
			}
		}
	}
}

// compareMaps
// It compares if two maps have the same content and return bool value. 
func compareMaps(tempMap map[string]string, fileContent map[string]string) bool {
	if len(tempMap) == len(fileContent) {
		for key, value := range tempMap {
			fileValue, ok := fileContent[key]
			if !ok {
				return false
			} else {
				if fileValue != value {
					return false
				}
			}
		}
	} else {
		return false
	}
	return true
}

// CreateSession
// This function generates univerally unique identifier for session and store it with name in map
func (s *Sessions) SetSession(name string, uuid string) {
	//add the name with uuid to map
	s.SessionsSyncLoc.Lock()
	s.SessionsMap[uuid] = name
	log.Debugf("Logged in user. Name: %s, uuid: %s", name, uuid)
	s.SessionsSyncLoc.Unlock()
}

// GetSession
// It checks if the key for the session exists in the map. It returns the bool value
// and the name if exist
func (s *Sessions) GetSession(key string) (string, bool) {
	var ok bool
	var name string
	s.SessionsSyncLoc.RLock()
	name, ok = s.SessionsMap[key]
	log.Debugf("Retreive session. Name: %s, uuid: %s", name, key)
	s.SessionsSyncLoc.RUnlock()
	return name, ok
}

// RemoveSession
// Removes the session from map
func (s *Sessions) RemoveSession(key string) {
	log.Debugf("Delete session. uuid: %s", key)
	delete(s.SessionsMap, key)
}
