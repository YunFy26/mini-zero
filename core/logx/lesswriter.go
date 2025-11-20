package logx

import "io"

type lessWriter struct {
	*limitedExecutor
	writer io.Writer
}

// newLessWriter creates a new lessWriter that limits the write frequency to the specified milliseconds.
func newLessWriter(writer io.Writer, milliseconds int) *lessWriter {
	return &lessWriter{
		limitedExecutor: newLimitedExecutor(milliseconds),
		writer:          writer,
	}
}

// Write writes data to the underlying writer with rate limiting.
func (lw *lessWriter) Write(p []byte) (n int, err error) {
	lw.logOrDiscard(func() {
		lw.writer.Write(p)
	})
	return len(p), nil
}
