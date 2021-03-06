// CSS 490
// Magdalena Grzmiel
// Assignments #3
// Copyright 2015 Magdalena Grzmiel
// sessionManager is resposnisible for managing the session.

package sessionManager

import (
	"bytes"
	log "github.com/cihub/seelog"
	"os/exec"
	"strings"
	"sync"
)

// Session structure keeps the map with the session keys and
// RWMutex object for the synchronization purpose
type Sessions struct {
	SessionsMap     map[string]string
	SessionsSyncLoc *sync.RWMutex
}

// New
// Creates new Session structure
func New() *Sessions {
	return &Sessions{SessionsMap: make(map[string]string), SessionsSyncLoc: new(sync.RWMutex)}
}

// CreateSession
// This function generates univerally unique identifier for session and store it with name in map
func (s *Sessions) CreateSession(name string) string {
	//generate uuid
	cmd := exec.Command("/usr/bin/uuidgen")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorf("Not able to generate uuid", err)
	}
	uuid := out.String()
	uuid = strings.Replace(uuid, "\n", "", 1)

	//add the name with uuid to map
	s.SessionsSyncLoc.Lock()
	s.SessionsMap[uuid] = name
	log.Debugf("Logged in user. Name: %s, uuid: %s", name, uuid)
	s.SessionsSyncLoc.Unlock()
	return uuid
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
