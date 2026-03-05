package bbolt

import (
	"fmt"
	"sync"
)

// MemData is an in-memory implementation of the Data interface.
type MemData struct {
	mu   sync.RWMutex
	data []byte
}

func NewMemData() *MemData {
	return &MemData{}
}

func (m *MemData) ReadAt(off int64, n int) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	end := int(off) + n
	if off < 0 || end > len(m.data) {
		return nil, fmt.Errorf("memdata: read [%d, %d) out of range [0, %d)", off, end, len(m.data))
	}
	return m.data[off:end], nil
}

func (m *MemData) WriteAt(b []byte, off int64) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	end := int(off) + len(b)
	if end > len(m.data) {
		m.data = append(m.data, make([]byte, end-len(m.data))...)
	}
	copy(m.data[off:], b)
	return len(b), nil
}

func (m *MemData) Size() (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return int64(len(m.data)), nil
}

func (m *MemData) Grow(sz int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if int64(len(m.data)) < sz {
		m.data = append(m.data, make([]byte, int(sz)-len(m.data))...)
	}
	return nil
}

func (m *MemData) Bytes() []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()
	dst := make([]byte, len(m.data))
	copy(dst, m.data)
	return dst
}

func (m *MemData) SetBytes(b []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make([]byte, len(b))
	copy(m.data, b)
}

func (m *MemData) Sync() error {
	return nil
}
