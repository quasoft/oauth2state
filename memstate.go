package oauth2state

import (
	"fmt"
	"sync"
)

// MemStateStore is an in-memory implementation of StateStorer that
// can be used safely by concurrent goroutines in a single server instance
type MemStateStore struct {
	states         map[string]string
	mutex          sync.RWMutex
	valueGenerator ValueGenerator
}

// NewMemStateStore creates a new memory store for state values
func NewMemStateStore() *MemStateStore {
	return &MemStateStore{
		states:         make(map[string]string),
		valueGenerator: CryptoValueGenerator{},
	}
}

// NewState creates a new random state value and associates the given URL with that value
func (s *MemStateStore) NewState(url string) (string, error) {
	state := s.valueGenerator.String()
	err := s.Add(state, url)
	return state, err
}

// Add a new state/url combination to the store
func (s *MemStateStore) Add(state, url string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if state == "" {
		return fmt.Errorf("State argument not provided")
	}

	s.states[state] = url

	return nil
}

// Contains checks if the given state value exists in the store.
func (s *MemStateStore) Contains(state string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if state == "" {
		return false, fmt.Errorf("State argument not provided")
	}

	_, ok := s.states[state]
	return ok, nil
}

// URL retrieves the URL that is associated with a given state value
func (s *MemStateStore) URL(state string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if state == "" {
		return "", fmt.Errorf("State argument not provided")
	}

	url, ok := s.states[state]
	if !ok {
		return "", fmt.Errorf("could not find state value %q in memory store", state)
	}
	return url, nil
}

// Delete the given state value from the store
func (s *MemStateStore) Delete(state string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if state == "" {
		return fmt.Errorf("State argument not provided")
	}

	delete(s.states, state)

	return nil
}
