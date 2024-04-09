package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type FileWithLine struct {
	Line int
	Path string
}

var wg = sync.WaitGroup{}

var mySet = map[string]FileWithLine{}
var rgx1 = regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)
var rgx2 = regexp.MustCompile(`vid/[a-f0-9]{12}`)

var rgx *regexp.Regexp

func init() {
	if os.Getenv("use_vid") == "1" || os.Getenv("use_vid") == "true" {
		rgx = rgx2
	} else {
		rgx = rgx1
	}
}

var mtx = sync.RWMutex{}

// ErrId XXX: "0a69f97b-b273-4d70-8061-f5eb85277d15",
// ErrId XXX: "0a69f97b-b273-4d70-8061-f5eb85277d15",

// "cmXXX:333333-b273-4d70-8061-f5eb85277d15",
// "cmXXX:333333-b273-4d70-8061-f5eb85277d15",

// "333333-b273-4d70-8061-f5eb85277d15:",
// "333333-b273-4d70-8061-f5eb85277d15:",

// "vid/f5eb85277d15",
// "vid/f5eb85277d15",

// "vid/f5eb85277d15:",
// "vid/f5eb85277d15:",

func traverseDir(d string) {

	// fmt.Println("traversing dir:", d)

	bn := filepath.Base(d)

	switch bn {
	case ".git",
		"tmp",
		"temp",
		".github",
		".idea",
		"pkg",
		"node_modules",
		".vscode":
		return
	}

	wg.Add(1)
	go func() {

		defer wg.Done()

		files, err := ioutil.ReadDir(d) // os.ReadDir() instead...

		if err != nil {
			log.Fatal("7819ea77-24ea-4c24-b11f-6d968e606bf5", err)
		}

		for _, f := range files {

			fullPath := filepath.Join(d, f.Name())
			// fmt.Println("we see file:", fullPath)

			if f.IsDir() {
				if strings.HasSuffix(fullPath, "/logs") {
					fmt.Println("[err-id-checker] skipping /logs path:", fullPath)
					continue
				}
				traverseDir(fullPath)
				continue
			}

			if strings.HasSuffix(fullPath, ".md") {
				fmt.Println("[err-id-checker] skipping .md file:", fullPath)
				continue
			}

			if strings.HasSuffix(fullPath, ".log") {
				fmt.Println("[err-id-checker] skipping .log file:", fullPath)
				continue
			}

			func(fullPath string) {

				file, err := os.Open(fullPath)

				if err != nil {
					log.Println(fmt.Sprintf("failed opening file: %v", err))
					return
				}

				defer func() {
					if err := file.Close(); err != nil {
						log.Println("12ee3ee5-1232-4ac0-9294-a376b764f9e0", err)
					}
				}()

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)

				i := 0
				for scanner.Scan() {
					i++
					var line = scanner.Text()

					var doThing = func(theUuid string) {

						if strings.HasPrefix(theUuid, "zz") {
							fmt.Println("[err-id-checker] skipping uuid with 'zz' in front")
							return
						}

						if strings.HasPrefix(theUuid, "cm:") {
							theUuid = theUuid[3:]
						}

						if strings.HasSuffix(theUuid, ":") {
							theUuid = theUuid[:len(theUuid)-1]
						}

						mtx.RLock()
						v, ok := mySet[theUuid]
						mtx.RUnlock()

						if ok {
							fmt.Println("the set already has this uuid:", theUuid, "at: ", fmt.Sprintf("%s:%v", v.Path, v.Line))
							fmt.Println("the current dupe is located at:", fmt.Sprintf("%s:%v", fullPath, i))
							fmt.Println("new uuid 1:", strings.ToLower(uuid.New().String()))
							fmt.Println("new uuid 2:", strings.ToLower(uuid.New().String()))
							fmt.Println("new uuid 3:", strings.ToLower("vid/"+(uuid.New().String()[24:36])))
							fmt.Println("new uuid 4:", strings.ToLower("vid/"+(uuid.New().String()[24:36])))
							log.Fatal("We found a dupe.")
						} else {
							mtx.Lock()
							mySet[theUuid] = FileWithLine{Line: i, Path: fullPath}
							mtx.Unlock()
						}
					}

					captured := rgx.FindStringSubmatch(line)

					if len(captured) > 1 {
						// fmt.Println("capture group:", captured)
						log.Fatal("eid/1115c905ebb8:", "strange capture length greater than 1")
					}

					if len(captured) < 1 {
						continue
					}

					var cap0 = captured[0]
					doThing(cap0)

				}
			}(fullPath)
		}

	}()

}

func main() {

	var dir = ""
	var next = false
	var lastValue = ""

	var args = os.Args

	if len(args) > 1 {
		lastValue = args[len(args)-1]
	}

	for _, s := range os.Args {
		if next {
			dir = s
			break
		}
		if s == "-d" {
			next = true
		}
		if s == "--dir" {
			next = true
		}
		if strings.HasPrefix(s, "--dir=") {
			dir = strings.Join(strings.Split(s, "=")[1:], "")
			break
		}
	}

	if dir == "" {
		dir = lastValue
		if dir == "" {
			fmt.Println("b4d000c6-63a4-400b-a5a6-b87d53a06056", "please provide a value for --dir, or -d, or --dir=x")
			os.Exit(1)
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("71c1589a-7384-4203-b575-63de85f7234c", err)
	}

	if strings.HasPrefix(dir, "--dir=") {
		dir = strings.Join(strings.Split(dir, "=")[1:], "")
	}

	if !filepath.IsAbs(dir) {
		dir = filepath.Join(pwd, "/", dir)
	}

	fmt.Println("the dir we are searching:", dir)

	if v, err := os.Stat(dir); true {
		if err != nil {
			fmt.Println("7be3285a-31b4-461d-b01c-721d7637e0ca", err)
			fmt.Println("Could not find directory specified by last/final arg or --dir")
			os.Exit(1)
		}
		if !v.IsDir() {
			fmt.Println("49d6eb7e-4e70-4eac-a4bf-738bbf1ec167", "not configured to search one file")
			fmt.Println("Please pass a directory, not a file, as argument to --dir, or as last/final argument in the command list.")
			os.Exit(1)
		}
	}

	fmt.Println("traversing root dir:", dir)
	traverseDir(dir)
	fmt.Println("Main: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed")
}
