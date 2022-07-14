package main

import (
	"fmt"
	"io"
	"os"
)

type CounterWriter struct {
	initWriter io.Writer
	counter int64
};

func (icw *CounterWriter) Write (p []byte) (int, error) {
	icw.counter += int64(len(p))
	return icw.initWriter.Write(p)
}

func (icw *CounterWriter) String() string {
	return fmt.Sprintf("%d bytes was written.", icw.counter)
}

func CountingWriter (w io.Writer) (io.Writer, *int64) {
	wrappedWriter := CounterWriter{w, 0}
	return &wrappedWriter, &wrappedWriter.counter;
}


func main() {
	wrstdout, cntr := CountingWriter(os.Stdout)
	fmt.Fprintln(wrstdout, "Hello friend!")

	fmt.Printf("pointer count is: %d\n", *cntr)
	fmt.Println(wrstdout)

	fmt.Fprintln(wrstdout, "Hello friend! <again>")

	fmt.Printf("pointer count is: %d\n", *cntr)
	fmt.Println(wrstdout)

}
