package worker

import (
	"sync"
)

type WorkerMap struct {
	lock *sync.RWMutex
	bm   map[string]SimpleWorker
}

func NewWorkerMap() *WorkerMap {
	return &WorkerMap{
		lock: new(sync.RWMutex),
		bm:   make(map[string]SimpleWorker),
	}
}

//Get from maps return the k's value
func (m *WorkerMap) Get(k string) SimpleWorker {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

//Size Get Size
func (m *WorkerMap) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.bm)
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *WorkerMap) Set(k string, v SimpleWorker) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Check Returns true if k is exist in the map.
func (m *WorkerMap) Check(k string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

func (m *WorkerMap) Delete(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

func (m *WorkerMap) Keys() []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var keys []string
	if len(m.bm) == 0 {
		return nil
	}
	for key, _ := range m.bm {
		keys = append(keys, key)
	}
	return keys
}

func (m *WorkerMap) Workers() []SimpleWorker {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var workers []SimpleWorker
	if len(m.bm) == 0 {
		return nil
	}
	for _, worker := range m.bm {
		workers = append(workers, worker)
	}
	return workers
}
