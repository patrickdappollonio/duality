// prefixer_test.go

package prefixer

import (
	"bytes"
	"testing"
)

// Test handling of input where the last line doesn't end with a newline.
func TestPrefixWriter_MissingFinalNewline(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	prefix := "[PREFIX] "
	pw := NewPrefixWriter(&buf, prefix)

	input := "Line 1\nLine 2\nLine 3"
	expectedOutput := prefix + "Line 1\n" +
		prefix + "Line 2\n" +
		prefix + "Line 3"

	n, err := pw.Write([]byte(input))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if n != len(input) {
		t.Fatalf("Expected to write %d bytes, but wrote %d bytes", len(input), n)
	}

	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%q\nBut got:\n%q", expectedOutput, buf.String())
	}
}
