package main

import (
	"sync"
)

var (
	visitedMap     = make(map[string]bool)
	visitedMapLock sync.Mutex
)

func visitUser(userID string) bool {
	visitedMapLock.Lock()
	defer visitedMapLock.Unlock()

	if visitedMap[userID] {
		return false
	}

	visitedMap[userID] = true
	return true
}
