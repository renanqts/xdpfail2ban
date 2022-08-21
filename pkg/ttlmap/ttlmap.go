package ttlmap

import (
	"sync"
	"time"
)

// TTLMap implements a map where the values expire with TTLs
type TTLMap struct {
	entries    map[interface{}]interface{}
	eMutex     sync.Mutex
	schedule   map[interface{}]*time.Timer
	sMutex     sync.Mutex
	onRemoval  func(interface{}, interface{})
	defaultTTL time.Duration
}

// New will create a TTLMap with a callback function
func New(ttl time.Duration, callback func(key, value interface{})) *TTLMap {
	return &TTLMap{
		make(map[interface{}]interface{}),
		sync.Mutex{},
		make(map[interface{}]*time.Timer),
		sync.Mutex{},
		callback,
		ttl,
	}
}

// clearSchedule removes expired entries from the schedule
func (t *TTLMap) clearSchedule(key interface{}) {
	t.sMutex.Lock()
	defer t.sMutex.Unlock()
	delete(t.schedule, key)
}

func (t *TTLMap) addEntry(key, value interface{}) {
	t.eMutex.Lock()
	defer t.eMutex.Unlock()
	t.entries[key] = value
}

func (t *TTLMap) addWithTTL(key, value interface{}, ttl time.Duration) {
	t.sMutex.Lock()
	defer t.sMutex.Unlock()
	if t.schedule[key] != nil {
		// Reset the ttl for entries that exist already
		t.schedule[key].Reset(ttl)
	} else {
		// create a timer to monitor it, when expires, then remove
		// the object from the entries
		t.schedule[key] = time.NewTimer(ttl)
		go func() {
			if timer, found := t.schedule[key]; found {
				<-timer.C
			}
			t.Delete(key)
		}()
	}
	t.addEntry(key, value)
}

// Add adds an entry with a specified TTL value
func (t *TTLMap) Add(key, value interface{}) {
	t.addWithTTL(key, value, t.defaultTTL)
}

// Delete removes an entry from the entries
func (t *TTLMap) Delete(key interface{}) {
	t.eMutex.Lock()
	defer t.eMutex.Unlock()
	delete(t.entries, key)

	// delete from the schedule
	t.clearSchedule(key)
}

// Get returns a value from the map
func (t *TTLMap) Get(key interface{}) (value interface{}) {
	t.eMutex.Lock()
	defer t.eMutex.Unlock()

	if value, present := t.entries[key]; present {
		return value
	}

	return nil
}
