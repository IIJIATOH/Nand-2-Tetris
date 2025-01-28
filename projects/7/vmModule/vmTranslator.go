package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

// M - память
// A - регистр. Принимает значение @номер, Также является указателем ячейки памяти
// D - регистр. Дополнительно хранит значения.

var segmentMap = map[string]string{
	"constant": "@SP",
	"local":    "@LCL",
	"argument": "@ARG",
	"this":     "@THIS",
	"that":     "@THAT",
}
var pointerSegmentMap = map[string]int{
	"temp":     5,
	"local":    300,
	"argument": 400,
	"this":     3000,
	"that":     3010,
}
var setMemoryOnSPValue = "@SP\nM=M-1\nA=M\n"
var addSP = "@SP\nM=M+1\n"

var counter int64 = 0

func Translate(path string) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	checkCounter := &counter

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

	logicalMap := map[string]string{
		"add": setMemoryOnSPValue + "D=M\n" + setMemoryOnSPValue + "M=M+D\n" + addSP,
		"sub": setMemoryOnSPValue + "D=M\n" + setMemoryOnSPValue + "M=M-D\n" + addSP,
		"neg": setMemoryOnSPValue + "M=-M\n" + addSP,
		"eq":  "JEQ",
		"gt":  "JGT",
		"lt":  "JLT",
		"and": setMemoryOnSPValue + "D=M\n" + setMemoryOnSPValue + "M=D&M\n" + addSP,
		"or":  setMemoryOnSPValue + "D=M\n" + setMemoryOnSPValue + "M=D|M\n" + addSP,
		"not": setMemoryOnSPValue + "M=!M\n" + addSP,
	}
	stackMap := map[string]func(string, string) string{
		"push": push,
		"pop":  pop,
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimLine, "//") || trimLine == "" {
			continue
		}
		if processedLine, exists := logicalMap[line]; exists {
			comparisons := []string{"eq", "gt", "lt"}
			if slices.Contains(comparisons, line) {
				_, err := writer.WriteString(comparisson(logicalMap[line]) + "\n")
				fmt.Println(err)
			} else {
				_, err := writer.WriteString(processedLine + "\n")
				fmt.Println(err)
			}
		} else {
			splittedLine := strings.Split(line, " ")
			if fn, exists := stackMap[splittedLine[0]]; exists {
				processedLine := fn(splittedLine[1], splittedLine[2])
				fmt.Println(checkCounter)
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

func push(segmentLine, number string) string {
	var result string
	strNumber, _ := strconv.Atoi(number)

	if segmentLine == "constant" {
		result = fmt.Sprintf(
			"@%s\nD=A\n%s\nA=M\nM=D\n%s\nM=M+1",
			number,
			segmentMap[segmentLine],
			segmentMap[segmentLine],
		)
	} else {
		lineNumber := pointerSegmentMap[segmentLine] + strNumber
		result = fmt.Sprintf(
			"@%d\nD=A\n@%d\nM=D\n",
			strNumber,
			lineNumber,
		)
	}

	return result
}

func pop(segmentLine string, number string) string {
	convertedNumber, _ := strconv.Atoi(number)
	result := "@SP\nM=M-1\nA=M\nD=M\n@" + strconv.Itoa(pointerSegmentMap[segmentLine]+convertedNumber) + "\nM=D\n"
	fmt.Println(result)
	return result
}
func comparisson(JMPKey string) string {
	counter += 2
	jmp := fmt.Sprintf("D=M-D\n@END%d\nD;%s\n@SP\nA=M\nM=-1\n@END%d\n0;JMP\n(END%d)\n@SP\nA=M\nM=1\n(END%d)\n", counter, JMPKey, counter+1, counter, counter+1)
	return setMemoryOnSPValue + "D=M\n" + setMemoryOnSPValue + jmp + addSP
}
