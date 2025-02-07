package main

import (
	"fmt"
	"log"
	"os"
)

func Translate(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	parser := NewParser(file) // Передаём file как io.Reader
	if parser.hasMoreLines() {
		parser.advance()
		fmt.Println(parser.currentCommand)
		parser.hasMoreLines()
		parser.advance()
		fmt.Println(parser.currentCommand)
		fmt.Println(parser.commandType())
		parser.arg1()
		parser.arg2()
	}
}
