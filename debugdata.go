package bbolt

// DebugData wraps a Data implementation and zeroes buffers returned by
// ReadAt when their release function is called.  This catches
// use-after-release bugs: any code that reads from a buffer after
// releasing it will see zeros instead of stale data, making the
// problem immediately visible.
type DebugData struct {
	Data
}

func (d *DebugData) ReadAt(off int64, n int) ([]byte, func(), error) {
	buf, release, err := d.Data.ReadAt(off, n)
	if err != nil {
		return nil, nil, err
	}
	// Allocate a larger backing buffer with a poison pattern before and
	// after the requested data.  Any code that reads past the returned
	// slice (e.g. via unsafe pointer arithmetic into a stale page) will
	// hit 0x67 bytes instead of valid data, making the bug obvious.
	const poisonLen = 256
	full := make([]byte, poisonLen+n+poisonLen)
	for i := 0; i < poisonLen; i++ {
		full[i] = 0x67
	}
	copy(full[poisonLen:], buf)
	for i := poisonLen + n; i < len(full); i++ {
		full[i] = 0x67
	}
	result := full[poisonLen : poisonLen+n : poisonLen+n]
	return result, func() {
		release()
		go clear(full)
	}, nil
}
