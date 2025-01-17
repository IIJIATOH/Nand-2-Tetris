package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Translate(path string) {
	type Operation func(string)
	// set RAM[0] 256,   // stack pointer
	// set RAM[1] 300,   // base address of the local segment
	// set RAM[2] 400,   // base address of the argument segment
	// set RAM[3] 3000,  // base address of the this segment
	// set RAM[4] 3010,  // base address of the that segment
	// var Addresses = `
	// @256
	// D=A
	// @SP
	// M=D
	// @300
	// M=D
	// @LCL
	// D=A
	// @300
	// M=D
	// @ARG
	// D=A
	// @400
	// M=D
	// @THIS
	// D=A
	// @3000
	// M=D
	// @THAT
	// D=A
	// @3010
	// M=D
	// `
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

	// Создаем сканер
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)
	// writer := bufio.NewWriter(outputFile)

	// • add, sub , neg
	// • eq , gt , lt
	// • and, or , not
	logicalMap := map[string]string{
		"add":  "@SP\nA=M\nD=M\n@SP\nM=M-1\nM=M+D",
		"sub":  "test",
		"neg":  "test",
		"eq":   "test",
		"gt":   "test",
		"lt":   "test",
		"and":  "test",
		"or":   "test",
		"not":  "test",
		"push": "test",
		"pop":  "test",
	}
	stackMap := map[string]func(string, string) string{
		"push": push,
		"pop":  push,
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimLine, "//") || trimLine == "" {
			continue
		}
		if processedLine, exists := logicalMap[line]; exists {
			_, err := writer.WriteString(processedLine + "\n")
			fmt.Println(err)
		} else {
			splittedLine := strings.Split(line, " ")
			if fn, exists := stackMap[splittedLine[0]]; exists {
				processedLine := fn(splittedLine[1], splittedLine[2])
				_, err := writer.WriteString(processedLine + "\n")
				fmt.Println(err)
			}

			fmt.Println(splittedLine[1])
		}
	}
	if err := writer.Flush(); err != nil {
		log.Fatalf("failed to flush writer: %s", err)
	}
}

func push(segmentLine string, number string) string {
	segmentMap := map[string]string{
		"constant": "@SP",
		"local":    "@LCL",
		"argument": "@ARG",
		"this":     "@THIS",
		"that":     "@THAT",
	}
	result := fmt.Sprintf("%s\nA=M\nD=M\n@%s\nM=A\n%s\nM=M+1", segmentMap[segmentLine], number, segmentMap[segmentLine])
	fmt.Println(result)
	return result
}
