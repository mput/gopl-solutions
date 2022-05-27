package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

var NOWORKFILE = errors.New("No error file found")

func searchForWorkFile(p string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if p == "/" || p == home {
		return "", NOWORKFILE
	}
	dirFiles, err := ioutil.ReadDir(p)
	if err != nil {
		return "", fmt.Errorf("can't list dir %s :%s", p, err)
	}
	for _, f := range dirFiles {
		if f.Name() == "go.work" {
			return path.Join(p, f.Name()), nil
		}

	}
	return searchForWorkFile(path.Dir(p))
}

func getWorkFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return searchForWorkFile(cwd)
}

func exname(s string) (string, error) {
	const regs = "^(\\d{1,2})\\.(\\d{1,3})$"
	r := regexp.MustCompile(regs)
	matches := r.FindStringSubmatch(s)
	if len(matches) != 3 {
		return "", fmt.Errorf("Arg should be in format '%s'", regs)
	}
	chapter, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}
	exr, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ex-%02d-%03d", chapter, exr), nil
}

func RunInDir(path, c string, args ...string) ([]byte, error) {
	cmd := exec.Command(c, args...)
	cmd.Dir = path
	return cmd.Output()
}

func NewEx(workfile, exName string) error {
	baseDir := path.Dir(workfile)
	exDir := path.Join(baseDir, exName)
	if err := os.Mkdir(exDir, 0755); err != nil {
		return err
	}
	log.Printf("dir '%s' created", exDir)

	_, err := RunInDir(exDir, "go", "mod", "init", fmt.Sprintf("mput.me/%s", exName))
	if err != nil {
		return err
	}
	log.Printf("mod '%s' inited", exName)

	_, err2 := RunInDir(baseDir, "go", "work", "use", fmt.Sprintf("./%s", exName))
	if err2 != nil {
		return err
	}
	log.Printf("record for '%s' added to 'go.work'", exName)

	gofile := fmt.Sprint(exName, ".go")
	mainpkgf := path.Join(exDir, gofile)
	f, err := os.Create(mainpkgf)
	defer f.Close()

	fmt.Fprintf(f, `package main

func main() {

}`)
	log.Printf("main package created in '%s' file", gofile)
	return nil
}

func DeleteEx(workfile, exName string) error {
	baseDir := path.Dir(workfile)
	exDir := path.Join(baseDir, exName)

	// Delete record wrom go.work file.
	wfread, err := os.Open(workfile)
	if err != nil {
		return err
	}
	wfcont, err := io.ReadAll(wfread)
	if err != nil {
		return err
	}
	wfread.Close()
	excRegexp := regexp.MustCompile(fmt.Sprintf("\\s*\\./%s", exName))
	res := excRegexp.ReplaceAll(wfcont, []byte(""))

	wfwrite, err := os.Create(workfile)
	if err != nil {
		return err
	}
	fmt.Fprint(wfwrite, string(res))
	wfwrite.Close()
	log.Printf("'go.work' file is cleaned from '%s' record", exName)

	// Delete exercise dir.
	if _, err := os.Stat(exDir); err != nil {
		return err
	}
	if err := os.RemoveAll(exDir); err != nil {
		return err
	}
	log.Printf("dir %s deleted", exDir)

	return nil
}

func main() {
	log.SetFlags(0)
	delFlag := flag.Bool("d", false, "delete exersise")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalf("You should specify exercise name")
	}

	workFile, err := getWorkFile()
	if err == NOWORKFILE {
		log.Fatalf("No 'go.work' were found in parent dirs.")
	} else if err != nil {
		log.Fatal(err)
	}

	exname, err := exname(flag.Arg(0))
	if err != nil {
		log.Fatalf("Can't format exercise dir name: %s", err)
	}


	if *delFlag {
		if err := DeleteEx(workFile, exname); err != nil {
			log.Fatalf("Fail! Can't delete exercise: %s", err)
		}
		fmt.Printf("Success! Exercise '%s' deleted\n", exname)
		return
	}

	if err := NewEx(workFile, exname); err != nil {
		log.Fatalf("Fail! Can't init exercise: %s", err)
	}
	fmt.Printf("Success! Exercise '%s' inited\n", exname)

}
