package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
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

var wg sync.WaitGroup

var mySet = map[string]FileWithLine{}
var rgx = regexp.MustCompile("\\s+ErrId\\s*:\\s*\"(?P<uuid>[0-9a-z-]{30,38})\"")
var rgx1 = regexp.MustCompile(".*\"cm:(?P<uuid>[0-9a-z-]{30,66})\"")
var rgx2 = regexp.MustCompile(".*\"cm:(?P<uuid>[0-9a-z-]{30,66}):\"")
var rgx3 = regexp.MustCompile(".*\"(?P<uuid>[0-9a-z-]{30,66}):\"")
var mtx = sync.RWMutex{}

// ErrId XXX: "0a69f97b-b273-4d70-8061-f5eb85277d15",
// ErrId XXX: "0a69f97b-b273-4d70-8061-f5eb85277d15",

// "cmXXX:333333-b273-4d70-8061-f5eb85277d15",
// "cmXXX:333333-b273-4d70-8061-f5eb85277d15",

// "333333-b273-4d70-8061-f5eb85277d15:",
// "333333-b273-4d70-8061-f5eb85277d15:",

func traverseDir(d string) {

	//fmt.Println("traversing dir:", d)

	bn := filepath.Base(d)

	if bn == ".git" {
		return
	}

	if bn == "tmp" {
		return
	}

	if bn == "temp" {
		return
	}

	if bn == ".github" {
		return
	}

	if bn == ".idea" {
		return
	}

	if bn == "node_modules" {
		return
	}

	if bn == ".vscode" {
		return
	}

	wg.Add(1)
	go func() {

		defer wg.Done()

		files, err := ioutil.ReadDir(d)

		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {

			fullPath := filepath.Join(d, f.Name())
			//fmt.Println("we see file:", fullPath)

			if f.IsDir() {
				traverseDir(fullPath)
				continue
			}

			file, err := os.Open(fullPath)

			defer file.Close()

			if err != nil {
				log.Fatalf("failed opening file: %s", err)
			}

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			i := 0
			for scanner.Scan() {
				i++
				var line = scanner.Text()

				var doThing = func(theUuid string) {

					//s := strings.Split(line, "\"")
					//theUuid := s[len(s)-2]

					//fmt.Println("the uuid before:", theUuid)

					if strings.HasPrefix(theUuid, "cm:") {
						theUuid = theUuid[3:]
					}

					if strings.HasSuffix(theUuid, ":") {
						theUuid = theUuid[:len(theUuid)-1]
					}

					//fmt.Println("the uuid:", theUuid)
					mtx.RLock()
					v, ok := mySet[theUuid]
					mtx.RUnlock()

					if ok {
						fmt.Println("the set already has this uuid:", theUuid, "at: ", fmt.Sprintf("%s:%v", v.Path, v.Line))
						fmt.Println("the current dupe is located at:", fmt.Sprintf("%s:%v", fullPath, i))
						fmt.Println("new uuid 1:", strings.ToLower(uuid.New().String()))
						fmt.Println("new uuid 2:", strings.ToLower(uuid.New().String()))
						fmt.Println("new uuid 3:", strings.ToLower(uuid.New().String()))
						log.Fatal("We found a dupe.")
					} else {
						mtx.Lock()
						mySet[theUuid] = FileWithLine{Line: i, Path: fullPath}
						mtx.Unlock()
					}
				}

				if rgx.MatchString(line) {
					//fmt.Println("line matches:", line)
					captured := rgx.FindStringSubmatch(line)
					if len(captured) > 1 {
						//fmt.Println("capture group:", captured)
					}

					if len(captured) > 1 && captured[1] != "" {
						doThing(captured[1])
					} else {
						fmt.Println(aurora.Red("capture group did not have expected length:"), captured)
					}
					continue
				}

				if rgx1.MatchString(line) {
					captured := rgx1.FindStringSubmatch(line)

					if len(captured) > 1 {
						//fmt.Println("capture group:", captured)
					}

					if len(captured) > 1 && captured[1] != "" {
						doThing(captured[1])
					} else {
						fmt.Println(aurora.Red("capture group did not have expected length:"), captured)
					}
					continue
				}

				if rgx2.MatchString(line) {
					captured := rgx2.FindStringSubmatch(line)
					if len(captured) > 1 {
						//fmt.Println("capture group:", captured)
					}

					if len(captured) > 1 && captured[1] != "" {
						doThing(captured[1])
					} else {
						fmt.Println(aurora.Red("capture group did not have expected length:"), captured)
					}
					continue
				}

				if rgx3.MatchString(line) {
					captured := rgx3.FindStringSubmatch(line)

					if len(captured) > 1 {
						//fmt.Println("capture group:", captured)
					}
					if len(captured) > 1 && captured[1] != "" {
						doThing(captured[1])
					} else {
						fmt.Println(aurora.Red("capture group did not have expected length:"), captured)
					}
					continue
				}

			}

		}

	}()

}

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	traverseDir(dir)
	fmt.Println("Main: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed")
}
