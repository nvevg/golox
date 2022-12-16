package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nvevg/golox/scanner"
)

var hadError = false

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func run(code string) {
}

func runFile(fileName string) {
	content, err := os.ReadFile(fileName)
	check(err)
	run(string(content))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _, e := reader.ReadLine()
		check(e)
		run(string(line))

		if hadError {
			os.Exit(65)
		}

		hadError = false
	}
}

const Code = `12name var_with_underscores var a = 12.34424.112 "a string" // a comment`

func main() {
	reader := bufio.NewReader(strings.NewReader(Code))
	sc := scanner.NewScanner(*reader)

	for t, e := sc.Scan(); e != io.EOF; t, e = sc.Scan() {
		if e != nil {
			fmt.Printf("%v\n", e)
			continue
		}
		fmt.Println(*t)
	}
}
