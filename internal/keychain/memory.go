package keychain

import "sync"

// MemoryKeychain is an in-memory keychain for testing
type MemoryKeychain struct {
	mu     sync.RWMutex
	tokens map[string]string
}

func NewMemoryKeychain() *MemoryKeychain {
	return &MemoryKeychain{tokens: make(map[string]string)}
}

func (m *MemoryKeychain) Set(profile, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[profile] = token
	return nil
}

func (m *MemoryKeychain) Get(profile string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tokens[profile], nil
}

func (m *MemoryKeychain) Delete(profile string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, profile)
	return nil
}
