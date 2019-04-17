package handlers

import "sync"

type KeyedRWMutex struct {
	l map[string]*sync.RWMutex
	m sync.Mutex
}

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

func NewKeyedRWMutex() *KeyedRWMutex {
	return &KeyedRWMutex{
		l: make(map[string]*sync.RWMutex),
	}
}
