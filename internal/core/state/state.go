package state

import "sync"

type State struct {
	data map[string]struct{}
	mu   sync.RWMutex
}

func NewManager() *State {
	return &State{
		data: make(map[string]struct{}),
	}
}

func (m *State) Add(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = struct{}{}
}

func (m *State) Remove(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *State) GetAll() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}

	return keys
}
