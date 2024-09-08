package decowriter_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/goaux/decowriter"
)

func TestWriter(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			prefix   string
			suffix   string
			expected string
		}{
			{"Empty input", "", "PREFIX: ", " :SUFFIX", ""},
			{"Single line", "Hello, World!", "PREFIX: ", " :SUFFIX", "PREFIX: Hello, World!"},
			{"Multiple lines", "Line 1\nLine 2\nLine 3", "-> ", " <-", "-> Line 1 <-\n-> Line 2 <-\n-> Line 3"},
			{"Empty lines", "\n\nContent\n\n", "# ", " #", "#  #\n#  #\n# Content #\n#  #\n"},
			{"Empty lines", "\n\n\n\n", "# ", " #", "#  #\n#  #\n#  #\n#  #\n"},
			{"No newline at end", "Text without newline", "> ", " <", "> Text without newline"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				type Buf interface {
					Write([]byte) (int, error)
					String() string
				}
				for _, buf := range []Buf{&bytes.Buffer{}, &strings.Builder{}} {
					w := decowriter.New(buf, []byte(tt.prefix), []byte(tt.suffix))

					n, err := io.WriteString(w, tt.input)
					if err != nil {
						t.Fatalf("Write error: %v", err)
					}
					if n != len(tt.input) {
						t.Errorf("Write returned %d, want %d", n, len(tt.input))
					}

					if got := buf.String(); got != tt.expected {
						t.Errorf("Got %q, want %q", got, tt.expected)
					}

					if w.Written() != int64(len(tt.expected)) {
						t.Errorf("Written() returned %d, want %d", w.Written(), len(tt.expected))
					}
				}
			})
		}
	})

	t.Run("", func(t *testing.T) {
		w := decowriter.New(&limitWriter{limit: 5}, []byte("123456789"), []byte("1234567890"))
		n, err := w.Write([]byte("ABCDEFG"))
		if err == nil {
			t.Error("should be error")
		}
		if n != 0 {
			t.Errorf("%d should be 0", n)
		}
	})

	t.Run("", func(t *testing.T) {
		w := decowriter.New(&limitWriter{limit: 5}, []byte("123"), []byte("1234567890"))
		n, err := w.Write([]byte("ABCDEFG"))
		if err == nil {
			t.Error("should be error")
		}
		if n != 2 {
			t.Errorf("%d should be 2", n)
		}
	})

	t.Run("", func(t *testing.T) {
		w := decowriter.New(&limitWriter{limit: 10}, []byte("123"), []byte("1234567890"))
		n, err := w.Write([]byte("ABCD\nFGHI\n"))
		if err == nil {
			t.Error("should be error")
		}
		if n != 5 {
			t.Errorf("%d should be 5", n)
		}
	})

	t.Run("", func(t *testing.T) {
		w := decowriter.New(&limitWriter{limit: 8}, []byte("12345"), []byte("1234"))
		n, err := w.Write([]byte("ABCDE\nFGHIJ\n"))
		if err == nil {
			t.Error("should be error")
		}
		if n != 3 {
			t.Errorf("%d should be 3", n)
		}
	})
}

type limitWriter struct {
	limit int
}

func (w *limitWriter) Write(p []byte) (int, error) {
	x := len(p)
	if w.limit < x {
		n := w.limit
		w.limit = 0
		return n, errors.New("limitWriter")
	}
	w.limit -= x
	return x, nil
}

func TestWriter_Flush(t *testing.T) {
	wf := &writerFlush{}
	dw := decowriter.New(wf, nil, nil)
	if _, err := dw.Write([]byte("hello\nworld")); err != nil {
		t.Error(err)
	}
	if _, err := dw.Write([]byte("hello\nworld\n")); err != nil {
		t.Error(err)
	}
	dw.Flush()
	want := 3
	if wf.flush != want {
		t.Errorf("Flush called %d times, want %d times", wf.flush, want)
	}
}

type writerFlush struct {
	flush int
}

func (w *writerFlush) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *writerFlush) Flush() error {
	w.flush++
	return nil
}
