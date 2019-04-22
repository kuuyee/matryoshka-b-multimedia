package handlers

import "sync"

// KeyedRWMutex is a set of RWMutex indexed by string key
type KeyedRWMutex struct {
	l map[string]*sync.RWMutex
	m sync.Mutex
}

// GetMutex returns the RWMutex with the given key
func (l *KeyedRWMutex) GetMutex(key string) *sync.RWMutex {
	l.m.Lock()
	defer l.m.Unlock()
	m, ok := l.l[key]
	if !ok {
		lock := new(sync.RWMutex)
		l.l[key] = lock
		return lock
	}
	return m
}

// NewKeyedRWMutex creates a new RWMutex
func NewKeyedRWMutex() *KeyedRWMutex {
	return &KeyedRWMutex{
		l: make(map[string]*sync.RWMutex),
	}
}
