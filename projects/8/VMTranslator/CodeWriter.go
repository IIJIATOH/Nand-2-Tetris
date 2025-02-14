package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
)

var setMemoryOnSPValue = "@SP\nM=M-1\nA=M\n"
var addSP = "@SP\nM=M+1\n"
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

type CodeWriter struct {
	Writer  *bufio.Writer
	counter int
}

func NewWriter(file io.Writer) *CodeWriter {
	return &CodeWriter{
		Writer: bufio.NewWriter(file),
	}
}
func (cw *CodeWriter) comparisson(JMPKey string) string {
	cw.counter += 2
	jmp := fmt.Sprintf("D=M-D\n@END%d\nD;%s\n@SP\nA=M\nM=-1\n@END%d\n0;JMP\n(END%d)\n@SP\nA=M\nM=1\n(END%d)\n", cw.counter, JMPKey, cw.counter+1, cw.counter, cw.counter+1)
	return setMemoryOnSPValue + "D=M\n" + setMemoryOnSPValue + jmp + addSP
}
func (cw *CodeWriter) writeArithmetic(command string) {
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

	processedLine, _ := logicalMap[command]
	comparisons := []string{"eq", "gt", "lt"}
	if slices.Contains(comparisons, command) {
		_, err := cw.Writer.WriteString(cw.comparisson(logicalMap[command]) + "\n")
		fmt.Println(err)
	} else {
		_, err := cw.Writer.WriteString(processedLine + "\n")
		fmt.Println(err)
	}

}

func (cw *CodeWriter) writePushPop(command Commands, segment string, index int) {
	splittedLine := strings.Split(segment, " ")
	if command == C_PUSH {
		cw.Writer.WriteString(cw.push(splittedLine[1], splittedLine[2]))
	} else {
		cw.Writer.WriteString(cw.pop(splittedLine[1], splittedLine[2]))
	}
}

func (cw *CodeWriter) writeLabel(label string) {
	cw.Writer.WriteString(fmt.Sprintf("(%s)", label))
}

func (cw *CodeWriter) writeGoto(label string) {
	cw.Writer.WriteString(fmt.Sprintf("@%s\n0;JMP", label))
}

func (cw *CodeWriter) writeIf(label string) {
	cw.Writer.WriteString(fmt.Sprintf("@SP\nA=M\nD=M\n@%s\nD;JGE", label))
}

func (cw *CodeWriter) push(segmentLine, number string) string {
	var result string
	strNumber, _ := strconv.Atoi(number)

	if segmentLine == "constant" {
		result = fmt.Sprintf(
			"@%s\nD=A\n%s\nA=M\nM=D\n%s\nM=M+1\n",
			number,
			segmentMap[segmentLine],
			segmentMap[segmentLine],
		)
	} else {
		lineNumber := pointerSegmentMap[segmentLine] + strNumber
		result = fmt.Sprintf(
			"@%d\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n",
			lineNumber,
		)
	}

	return result
}

func (cw *CodeWriter) writeCall(functionName string, nArgs string) {
	numberArgs, _ := strconv.Atoi(nArgs)
	// Перекидываем аргументы
	for i := 0; i < numberArgs; i++ {
		cw.push("argument", strconv.Itoa(i))
	}
	// Делаем прыжок на функцию
	cw.Writer.WriteString(fmt.Sprintf("@%s\n0;JMP", functionName))
}

func (cw *CodeWriter) pop(segmentLine string, number string) string {
	convertedNumber, _ := strconv.Atoi(number)
	result := "@SP\nM=M-1\nA=M\nD=M\n@" + strconv.Itoa(pointerSegmentMap[segmentLine]+convertedNumber) + "\nM=D\n"
	fmt.Println(result)
	return result
}

func (cw *CodeWriter) close() {
	if err := cw.Writer.Flush(); err != nil {
		log.Fatalf("failed to flush writer: %s", err)
	}
}
