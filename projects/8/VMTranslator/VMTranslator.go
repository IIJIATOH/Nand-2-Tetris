package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func Translate(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Создаем новый файл для записи
	outputFile, err := os.Create("output.asm")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	parser := NewParser(file) // Передаём file как io.Reader
	writer := NewWriter(outputFile)
	for parser.hasMoreLines() {
		parser.advance()
		trimLine := strings.TrimSpace(parser.currentCommand)
		if strings.HasPrefix(trimLine, "//") || trimLine == "" {
			continue
		}
		lineSlices := strings.Split(trimLine, " ")
		switch parser.currentCommandType {
		case C_ARITHMETIC:
			writer.writeArithmetic(parser.currentCommand)
		case C_POP:
			writer.writePushPop(parser.currentCommandType, parser.currentCommand, 0)
		case C_PUSH:
			writer.writePushPop(parser.currentCommandType, parser.currentCommand, 0)
		case C_LABEL:
			writer.writeLabel(lineSlices[1])
		case C_GOTO:
			writer.writeGoto(lineSlices[1])
		case C_IF:
			writer.writeIf(lineSlices[1])
		case C_CALL:
			writer.writeCall(lineSlices[0], lineSlices[1])

		default:
			continue
		}
	}
	writer.close()
}
