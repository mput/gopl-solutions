package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var port = flag.Int("port", 8080, "port number to run server at")
var c = flag.Bool("client", false, "whether to run in client mode")

func main() {
	flag.Parse()

	if !*c {
		runServer()
	} else {
		runClients()
	}
}

func runServer() {
	addr := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start listen at %s", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("Connection Accepted from %s", conn.RemoteAddr())
		go handleConn(conn)

	}

}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, fmt.Sprintln(time.Now().Format("15:04:05.00000")))
		log.Printf("Send ts to %s\n", c.RemoteAddr())
		if err != nil {
			log.Print(err)
			return
		}
		time.Sleep(1 * time.Second)
	}

}

const errMsg = "Specify servers in the format <name>=<address>"

type argument struct{Name, Addr string}

func runClients() {
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal(errMsg)
	}
	var argsP  []argument
	for i, arg := range args {
		argParsed := strings.Split(arg, "=")
		if len(argParsed) != 2 {
			log.Fatalf("(arg %d) Expected: %s, but '%s'", i+1, errMsg, arg)
			return
		}
		argsP = append(argsP, argument{argParsed[0], argParsed[1]})
	}
	width := 18
	printHeader(argsP, width)
	for i, argP := range argsP {
		go runClient(argP.Name, argP.Addr, i, len(argsP), width)
	}
	time.Sleep(1 * time.Hour)

}


func printHeader (args []argument, w int) {
	format := fmt.Sprintf("|%%%ds", w)
	for _, elm := range args {
		fmt.Printf(format, elm.Name)
	}
	fmt.Printf("|\n")
	for range args {
		fmt.Printf(format, strings.Repeat("-", w))
	}
	fmt.Printf("|\n")

}

func runClient(name, address string, idx, total, w int) {
	conn, err := net.Dial("tcp", address)
	formatPart := fmt.Sprintf("|%%%ds", w)
	formatFull := ""
	for i := 0; i < total; i++ {
		if i != idx {
			formatFull += "|" + strings.Repeat(" ", w)
		} else {
			formatFull += formatPart
		}
	}
	formatFull += "|\n"
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	mustCopyLinesf(formatFull, os.Stdout, conn)

}

func mustCopyLinesf(format string, dst io.Writer, src io.Reader) {
	linesScanner := bufio.NewScanner(src)
	for linesScanner.Scan() {
		line := linesScanner.Text()
		_, err := fmt.Fprintf(dst, format, line)
		if err != nil {
			log.Printf("Error while printing data: %s", err)
		}
	}
	if err := linesScanner.Err(); err != nil {
		log.Printf("Error while reading data: %s", err)
	}
}
