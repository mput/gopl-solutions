package main

import (
	"flag"
	"fmt"
	"mput.me/gopl/tempflags"
)


var temp = tempflags.CelsiusFlag("temp", 36.6, "body temperature")

func main() {
	flag.Parse()
	fmt.Printf("Setted temperature is %s\n", *temp)
}
