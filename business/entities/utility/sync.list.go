package utility

import (
	"sync"
)

type SyncList struct {
	sync.RWMutex
	data []interface{}
}

func (sl *SyncList) Get(index int) interface{} {
	sl.RLock()
	defer sl.RUnlock()
	return sl.data[index]
}

// this is not required for our usecase, but can be required for a list
func (sl *SyncList) insert(index int, val interface{}) {
	sl.Lock()
	defer sl.Unlock()
	if len(sl.data) <= index {
		newdata := make([]interface{}, index+len(sl.data))
		copy(newdata, sl.data)
		sl.data = newdata
	}
	sl.data[index] = val
}

func (sl *SyncList) Append(val interface{}) {
	sl.Lock()
	defer sl.Unlock()
	sl.data = append(sl.data, val)
}

// this is required for activity usecase
func (sl *SyncList) GetAllEntryList() []interface{} {
	sl.RLock()
	defer sl.RUnlock()
	data := make([]interface{}, len(sl.data))
	copy(data, sl.data)
	sl.data = nil
	return data
}
