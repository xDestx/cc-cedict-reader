package main

import (
	"bufio"
	"io"
	"log"
	"os"

	cccedictparser "github.com/xDestx/cc-cedict-reader"
)

func main() {
	var input io.Reader
	if len(os.Args) == 2 {
		// file
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
			return
		}
		defer f.Close()
		input = f
	} else if len(os.Args) == 1 {
		input = os.Stdin
	} else {
		log.Fatalf("format: <cmd> <optional file path>\nEx: cc-cedict-parse\nEx: cc-cedict-parse path/to/my/file")
		return
	}

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		l := scanner.Text()

		ci, err := cccedictparser.ParseLine(l)

		if err == nil {
			os.Stdout.WriteString(ci.String() + "\n")
		} else {
			os.Stdout.WriteString(err.Error())
		}
	}
}
