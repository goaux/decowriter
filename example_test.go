package decowriter_test

import (
	"bufio"
	"fmt"
	"os"

	"github.com/goaux/decowriter"
)

func Example() {
	w := decowriter.New(bufio.NewWriter(os.Stdout), []byte(">>"), []byte("<<"))
	fmt.Fprint(w, "hello")
	fmt.Fprintln(w, " world")
	fmt.Fprintln(w, "hello")
	fmt.Fprint(w, "world")
	fmt.Printf("\ntotal: %d\n", w.Written())
	// output:
	// >>hello world<<
	// >>hello<<
	// >>world
	// total: 33
}

func Example_readme() {
	w := decowriter.New(bufio.NewWriter(os.Stdout), []byte(">>"), []byte("<<"))
	fmt.Fprintln(w, "This is a log message")
	fmt.Fprintln(w, "Another log message")
	// Output:
	// >>This is a log message<<
	// >>Another log message<<
}
