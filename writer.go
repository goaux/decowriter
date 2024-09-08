// Package decowriter provides a writer that add a prefix and a suffix to each line.
//
// This package is the successor of github.com/goaux/prefixwriter.
// Compared to prefixwriter, it is the caller's responsibility to use bufio with decowriter.
// decowriter uses the underlying io.Writer directly.
package decowriter

import (
	"bufio"
	"bytes"
	"io"
	"slices"
)

// Writer is the writer that add a prefix and a suffix to each line.
// It implements the io.Writer interface.
//
// Writer defines a line as a sequence of zero or more non-'\n' bytes followed by a '\n'.
// A prefix is written before one or more bytes of a line.
// A suffix is written just before the trailing '\n'.
// A string without a trailing '\n' is not recognized as a line, so no suffix is written.
// This allows you to split a single line write across two Write calls.
type Writer struct {
	w io.Writer

	flush func() error

	prefix, suffix []byte

	head    bool
	written int64
}

// Ensure Writer implements io.Writer interface.
var _ io.Writer = (*Writer)(nil)

// New creates a new Writer that wraps the given io.Writer.
func New(w io.Writer, prefix, suffix []byte) *Writer {
	return &Writer{
		w:      w,
		flush:  getFlush(w),
		prefix: slices.Clone(prefix),
		suffix: append(append(make([]byte, 0, len(suffix)+1), suffix...), '\n'),
		head:   true,
	}
}

func getFlush(w io.Writer) func() error {
	if i, ok := w.(flusher); ok {
		return i.Flush
	}
	return nop
}

func nop() error { return nil }

type flusher interface {
	Flush() error
}

var _ flusher = (*bufio.Writer)(nil)

// Written returns the total number of bytes written, including prefixes.
func (w *Writer) Written() int64 {
	return w.written
}

// Write implements the io.Writer interface. It writes the given byte slice to the
// underlying writer, adding the prefix at the beginning of each line,
// adding the suffix at the end of each line.
//
// The returned int n represents the number of bytes from the input slice p that were
// processed, not including any added prefixes nor suffixes.
// This means that n <= len(p), even though the actual number of bytes written to the
// underlying writer may be larger due to the added prefixes and suffixes.
//
// If p contains no data (len(p) == 0), Write will not perform any operation and will
// return n = 0 and a nil error.
//
// An error is returned if the underlying writer returns an error, or if the Write
// operation cannot be completed fully.
//
// If one or more bytes are written to the underlying writer and the processing
// is completed normally, and if the underlying writer has `Flush() error` method,
// Flush is called once at the end.
func (w *Writer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	defer func() {
		if err == nil {
			err = w.flush()
		}
	}()

	for len(p) > 0 {
		if w.head {
			nn, err := w.w.Write(w.prefix)
			w.written += int64(nn)
			if err != nil {
				return n, err
			}
			w.head = false
		}

		i := bytes.IndexByte(p, '\n')
		if i == -1 {
			nn, err := w.w.Write(p)
			n += nn
			w.written += int64(nn)
			return n, err
		}

		nn, err := w.w.Write(p[:i])
		n += nn
		w.written += int64(nn)
		if err != nil {
			return n, err
		}

		nn, err = w.w.Write(w.suffix)
		n++ // for '\n'
		w.written += int64(nn)
		if err != nil {
			return n, err
		}

		p = p[i+1:]
		w.head = true
	}

	return n, nil
}

// Flush calls the Flush method of underlying writer if it has, otherwise returns nil.
func (w *Writer) Flush() error {
	return w.flush()
}
