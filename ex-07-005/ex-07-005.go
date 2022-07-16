package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type lr struct {
	reader io.Reader
	left int64
}

func (ilr *lr) Read(p []byte) (n int, e error) {
	if ilr.left == 0 {
		return 0, io.EOF
	}
	if ilr.left < int64(len(p)) {
		p = p[:ilr.left]
	}
	n, e = ilr.reader.Read(p)
	ilr.left -= int64(n);
	return
}


func LimitReader(r io.Reader, l int64) io.Reader {
	return &lr{r, l}
}

func main() {

	r, e := io.ReadAll(LimitReader(strings.NewReader("I need somebody. Help me!(("), 26))
	if e != nil {
		log.Fatal(e)
	}
	fmt.Printf("Result: (%d) %s\n", len(r), string(r))
}
