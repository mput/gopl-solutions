package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func startFTPServer(addr, basedir string) {
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
		ctx := new(sessionContext)
		ctx.cmdConn = conn
		ctx.baseDir = basedir
		go handleSession(ctx)

	}

}

type sessionContext struct {
	cmdConn           net.Conn
	baseDir           string
	currentDir        string
	requestedDataAddr string
}

func handleSession(ctx *sessionContext) {
	respond(ctx, "220")
	cmdScanner := bufio.NewScanner(ctx.cmdConn)
	for cmdScanner.Scan() {
		cmd := strings.TrimSpace(cmdScanner.Text())
		finish := handleCommand(ctx, cmd)
		if finish {
			break
		}
	}
	if err := cmdScanner.Err(); err != nil {
		log.Printf("Error while reading data: %s", err)
	}
}

func handleCommand(ctx *sessionContext, c string) (isFinish bool) {
	log.Printf("<--- %s", c)
	cmd, err := parseCmd(c)
	if err != nil {
		log.Printf("Command parsing error: %s", err)
		respond(ctx, "500 Syntax error")
		return false
	}
	switch cmd.Name {
	case "USER":
		respond(ctx, "230")
	case "SYST":
		respond(ctx, "215 Unix Type: L8")
	case "EPRT":
		setDataAddrFromEPRT(ctx, cmd)
	case "PORT":
		setDataAddrFromPORT(ctx, cmd)
	case "LIST":
		list(ctx)
	case "CWD":
		cwd(ctx, cmd)
	case "CDUP":
		cdup(ctx)
	case "RETR":
		get(ctx, cmd)
	case "QUIT":
		quit(ctx)
		return true
	default:
		respond(ctx, "502 Command not implemented")

	}
	return false
}

func currentPath(basedir, curdir string) string {
	return fmt.Sprintf("%s%s", basedir, curdir)
}

func sessionCurrentPath(ctx *sessionContext) string {
	return currentPath(ctx.baseDir, ctx.currentDir)
}

func quit(ctx *sessionContext) {
	err := ctx.cmdConn.Close()
	if err != nil {
		log.Print(err)
	}
	log.Printf("Session closed for %s", ctx.cmdConn.RemoteAddr())
}

func getFileList(ctx *sessionContext) (io.Reader, error) {
	cmd := exec.Command("ls", "-l", sessionCurrentPath(ctx))
	flist, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	res := []string{}
	lines := bufio.NewScanner(bytes.NewReader(flist))
	for lines.Scan() {
		line := lines.Text()
		if !strings.HasPrefix(line, "total") {
			res = append(res, line)
		}
	}

	res = append(res, "")
	if err := lines.Err(); err != nil {
		log.Printf("Error while reading data: %s", err)
	}
	if err != nil {
		return nil, err
	}
	return strings.NewReader(strings.Join(res, "\r\n")), nil
}

func list(ctx *sessionContext) {
	fileList, err := getFileList(ctx)
	if err != nil {
		log.Print(err)
		respond(ctx, "550 Requested action not taken")
		return
	}
	sendData(ctx, fileList)
}

func get(ctx *sessionContext, cmd Command) {
	if cmd.Arg == "" {
		respond(ctx, "502 Syntax error in parameters or argument")
		return
	}
	filePath := fmt.Sprintf("%s/%s", sessionCurrentPath(ctx), cmd.Arg)
	file, err := os.Open(filePath)
	if err != nil {
		log.Print(err)
		respond(ctx, "550")
		return
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if errors.Is(err, os.ErrNotExist) {
		respond(ctx, fmt.Sprintf("550 %s not found", cmd.Arg))
		return
	}
	if err != nil {
		log.Printf("Error at ReadFIle attempt to %s: %s", cmd.Arg, err)
		respond(ctx, "451 Requested action aborted: local error in processing")
		return
	}
	if fileStat.IsDir() {
		respond(ctx, fmt.Sprintf("550 %s is a file, not directory", cmd.Arg))
		return
	}
	sendData(ctx, file)
}

func dirExists(s string) bool {
	return true
}

func cwd(ctx *sessionContext, cmd Command) {
	if cmd.Arg == "" {
		respond(ctx, "502 Syntax error in parameters or argument")
		return
	}
	newCurrentDir := fmt.Sprintf("%s/%s", ctx.currentDir, cmd.Arg)
	fs, err := os.Stat(currentPath(ctx.baseDir, newCurrentDir))
	if errors.Is(err, os.ErrNotExist) {
		respond(ctx, fmt.Sprintf("550 %s not found", cmd.Arg))
		return
	}
	if err != nil {
		log.Printf("Error at CWD attempt to %s: %s", cmd.Arg, err)
		respond(ctx, "451 Requested action aborted: local error in processing")
		return
	}
	if fs.IsDir() {
		ctx.currentDir = newCurrentDir
		respond(ctx, "200")
		return
	} else {
		respond(ctx, fmt.Sprintf("550 %s is a file, not directory", cmd.Arg))
	}

}

func cdup(ctx *sessionContext) {
	ctx.currentDir = path.Dir(ctx.currentDir)
	respond(ctx, "200")
}

func setDataAddrFromEPRT(ctx *sessionContext, cmd Command) {
	a := make([]string, 0)
	for _, el := range strings.Split(cmd.Arg, cmd.Arg[0:1]) {
		if el != "" {
			a = append(a, el)
		}
	}
	if len(a) != 3 {
		log.Printf("wrong remote address argument: %s", cmd.Arg)
		respond(ctx, "502 Syntax error in parameters or argument")
		return
	}
	remAddr := fmt.Sprintf("%s:%s", a[1], a[2])
	ctx.requestedDataAddr = remAddr
	respond(ctx, "200")
}


func setDataAddrFromPORT(ctx *sessionContext, cmd Command) {
	a := strings.Split(cmd.Arg, ",")
	if len(a) != 6 {
		log.Printf("wrong remote address argument %s", cmd.Arg)
		respond(ctx, "502 Syntax error in parameters or argument")
	}
	p1, err1 := strconv.Atoi(a[4])
	p2, err2 := strconv.Atoi(a[5])
	if err1 != nil || err2 != nil {
		log.Printf("wrong remote address argument %s", cmd.Arg)
		respond(ctx, "502 Syntax error in parameters or argument")
		return
	}
	port := p1 * 256 + p2
	remAddr := fmt.Sprintf("%s.%s.%s.%s:%d", a[0], a[1], a[2], a[3], port)
	ctx.requestedDataAddr = remAddr
	respond(ctx, "200")
}

func sendData(ctx *sessionContext, r io.Reader) {
	if ctx.requestedDataAddr == "" {
		log.Print("Remote data address unknown")
		respond(ctx, "425 Can't open data connection")
		return
	}
	respond(ctx, "150 Opening data connection")
	conn, err := net.Dial("tcp", ctx.requestedDataAddr)
	if err != nil {
		log.Printf("can't open connection to %s: %s", ctx.requestedDataAddr, err)
		respond(ctx, "425 Can't open data connection")
		return
	}
	defer conn.Close()
	io.Copy(conn, r)
	respond(ctx, "226 Transfer complete")
}

type Command struct {
	Name string
	Arg  string
}

func parseCmd(c string) (Command, error) {
	var cmd Command
	a := strings.Fields(c)
	if len(a) < 1 {
		return cmd, fmt.Errorf("No command name present: %s", c)
	}
	cmd.Name = a[0]
	cmd.Arg = strings.Join(a[1:], " ")
	return cmd, nil
}

func respond(ctx *sessionContext, cmd string) error {
	log.Printf("---> %s", cmd)
	_, err := fmt.Fprintf(ctx.cmdConn, "%s\r\n", cmd)
	if err != nil {
		log.Printf("Response error: %s", err)
	}
	return err
}

func main() {
	startFTPServer("localhost:8020", ".")
}
