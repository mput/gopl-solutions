package main

import (
	"fmt"
	"io"
	"log"
)

type StringReader struct {
	s []byte
	p int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (sr *StringReader) Read(p []byte) (toSend int, e error) {
	if len(p) == 0 {
        return 0, nil
    }

	leftToSend := len(sr.s) - sr.p
	toSend = min(leftToSend, len(p))
	atEot := leftToSend <= toSend

	copy(p, sr.s[sr.p:])
	sr.p += toSend
	if atEot {
		e = io.EOF
	}
	return
}

func NewReader(s string) *StringReader {
	return &StringReader{[]byte(s), 0}
}

func main() {
	strReader := NewReader("I need somebody. Help me!")

	slc := make([]byte, 3)
	for {
		l, e := strReader.Read(slc)
		fmt.Printf("%d ->> '%s'\n", l, string(slc[0:l]))
		if e == io.EOF {
			break
		}

	}

	r, e := io.ReadAll(NewReader("I need somebody. Help me!(("))
	if e != nil {
		log.Fatal(e)
	}
	fmt.Printf("Result: (%d) %s", len(r), string(r))

}
