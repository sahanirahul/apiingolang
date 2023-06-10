package utility

import "sync"

type SyncMap struct {
	sync.RWMutex
	data map[string]interface{}
}

func (sm *SyncMap) Get(key string) interface{} {
	sm.RLock()
	defer sm.RUnlock()
	return sm.data[key]
}

func (sm *SyncMap) Put(key string, val interface{}) {
	sm.Lock()
	defer sm.Unlock()
	if sm.data == nil {
		sm.data = make(map[string]interface{})
	}
	sm.data[key] = val
}

func (sm *SyncMap) ContainsKey(key string) bool {
	sm.RLock()
	defer sm.RUnlock()
	if _, ok := sm.data[key]; !ok {
		return false
	}
	return true
}

// this is required for activity usecase
func (sm *SyncMap) GetAllEntry() map[string]interface{} {
	sm.RLock()
	defer sm.RUnlock()
	data := make(map[string]interface{})
	for key, val := range sm.data {
		data[key] = val
	}
	return data
}

// this is required for activity usecase
func (sm *SyncMap) PutIfNotPresent(key string, val interface{}) bool {
	sm.Lock()
	defer sm.Unlock()
	if sm.data == nil {
		sm.data = make(map[string]interface{})
	}
	if _, ok := sm.data[key]; ok {
		return false // key already presnt
	}
	sm.data[key] = val
	return true
}

// this is required for activity usecase
func (sm *SyncMap) GetAllEntryList() []interface{} {
	sm.RLock()
	defer sm.RUnlock()
	data := []interface{}{}
	for key, val := range sm.data {
		data = append(data, val)
		delete(sm.data, key)
	}
	return data
}
