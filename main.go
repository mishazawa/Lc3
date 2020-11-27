package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	runtime "github.com/mishazawa/Lc3/runtime"
	types "github.com/mishazawa/Lc3/runtime/types"
)

func main() {
	r := runtime.Boot()

	err := load(r)

	if err != nil {
		panic(err)
	}

	fmt.Println("\nHalt. Code:", r.Run())
}

var HELP_MESSAGE string = `
  help                     - print this message
  load "... path to image" - load assembler image
`

func load(r types.Runtime) error {
	var err error
	var mess string

	file := flag.String("exec", "", "Path to asm file")
	flag.Parse()

	if len(*file) != 0 {
		err, mess = loadFile(*file, r)
		if len(mess) != 0 {
			fmt.Println(mess)
			return nil
		}
		return err
	}

	COMMAND := regexp.MustCompile(`(load|help)\s?(.*)?$`)
	scanner := bufio.NewScanner(os.Stdin)

load_loop:
	for {
		fmt.Printf("Lc3 > ")
		scanner.Scan()

		groups := COMMAND.FindStringSubmatch(scanner.Text())

		if len(groups) == 0 {
			continue load_loop
		}

		switch groups[1] {
		case "help":
			fmt.Println(HELP_MESSAGE)
		case "load":
			err, mess = loadFile(groups[2], r)
			if len(mess) != 0 {
				fmt.Println(mess)
			} else {
				break load_loop
			}
		}
	}

	return err
}

func readImageToMemory(filename string, rt types.Runtime) error {
	file, err := os.Open(filename)

	if err != nil {
		return err
	}
	defer file.Close()

	return rt.Load(file)
}

func loadFile(path string, rt types.Runtime) (error, string) {
	if len(path) == 0 {
		return nil, "[Error] Enter path to file."
	}

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil, "[Error] File doesn't exist."
	} else if info.IsDir() {
		return nil, fmt.Sprintf("[Error] %s is directory.", path)
	} else {
		err := readImageToMemory(path, rt)
		return err, ""
	}
}
