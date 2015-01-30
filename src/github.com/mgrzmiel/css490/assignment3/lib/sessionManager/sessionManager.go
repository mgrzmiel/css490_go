package sessionManager

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type Sessions struct {
	SessionsMap     map[string]string
	SessionsSyncLoc *sync.RWMutex
}

func New() *Sessions {
	return &Sessions{SessionsMap: make(map[string]string), SessionsSyncLoc: new(sync.RWMutex)}
}

// generateUniqueId
// This function generates univerally unique identifier for cookie
func (s *Sessions) CreateSession(name string) string {
	cmd := exec.Command("/usr/bin/uuidgen")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	uuid := out.String()
	uuid = strings.Replace(uuid, "\n", "", 1)
	s.SessionsSyncLoc.Lock()
	s.SessionsMap[uuid] = name
	s.SessionsSyncLoc.Unlock()
	return uuid
}

// getNameAndCookie
// It checks if the cookie is set up and if the name for that cookie exists in map.
// Based on that, it sets up the correctlyLogIn variable.
func (s *Sessions) GetSession(key string) (string, bool) {
	var ok bool
	var name string
	s.SessionsSyncLoc.RLock()
	name, ok = s.SessionsMap[key]
	s.SessionsSyncLoc.RUnlock()
	return name, ok
}

func (s *Sessions) RemoveSession(key string) {
	delete(s.SessionsMap, key)
}
