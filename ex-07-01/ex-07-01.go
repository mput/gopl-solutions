package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type Wc struct {
	words int
	lines int
};

func (wc *Wc) Write(p []byte) (int, error) {
	scnw := bufio.NewScanner(bytes.NewReader(p))
	scnw.Split(bufio.ScanWords)
	for scnw.Scan() {
		*&wc.words++
	}

	scnl := bufio.NewScanner(bytes.NewReader(p))
	for scnl.Scan() {
		*&wc.lines++
	}

	return len(p), nil
}

func (wc *Wc) String() string {
	return fmt.Sprintf("Words count: %d\nLines count: %d", wc.words, wc.lines)
}


func main() {
	var wc Wc
	fmt.Fprint(&wc, "help me\n my friend here")

	fmt.Printf("%s", &wc)

}
