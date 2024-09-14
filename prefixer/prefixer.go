// prefixer.go

package prefixer

import (
	"bytes"
	"io"
	"sync"
)

// Package-level buffer pool using *bytes.Buffer.
var bufferPool = sync.Pool{
	New: func() interface{} {
		// Initialize a new bytes.Buffer with a predefined capacity.
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

// PrefixWriter wraps an io.Writer and prefixes each line with a specified string before writing.
type PrefixWriter struct {
	w           io.Writer
	prefixBytes []byte
}

// NewPrefixWriter creates a new PrefixWriter that writes to w with each line prefixed by prefix.
func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{
		w:           w,
		prefixBytes: []byte(prefix),
	}
}

// Write implements the io.Writer interface.
// It prefixes each line (complete or incomplete) and writes it to the underlying writer.
func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	start := 0

	for start < len(p) {
		// Find the next newline character
		idx := bytes.IndexByte(p[start:], '\n')
		var line []byte

		if idx == -1 {
			// No newline found; take the rest of the input
			line = p[start:]
			start = len(p)
		} else {
			// Include the newline character
			line = p[start : start+idx+1]
			start += idx + 1
		}

		// Get a bytes.Buffer from the pool
		outBuf := bufferPool.Get().(*bytes.Buffer)
		outBuf.Reset()

		// Build the prefixed line
		outBuf.Write(pw.prefixBytes)
		outBuf.Write(line)

		// Write the prefixed line to the underlying writer
		if _, err := pw.w.Write(outBuf.Bytes()); err != nil {
			bufferPool.Put(outBuf)
			return n, err
		}

		// Put the buffer back to the pool
		bufferPool.Put(outBuf)
	}

	return n, nil
}
