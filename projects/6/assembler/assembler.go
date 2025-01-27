package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func Translate(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	file2, err2 := os.Open(path)
	if err2 != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file2.Close()

	// Создаем новый файл для записи
	outputFile, err := os.Create("output.hack")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	// Создаем сканер
	scannerFirst := bufio.NewScanner(file)
	scannerSecond := bufio.NewScanner(file2)
	writer := bufio.NewWriter(outputFile)

	symbolMap := map[string]string{
		"R0":     "0",
		"R1":     "1",
		"R2":     "2",
		"R3":     "3",
		"R4":     "4",
		"R5":     "5",
		"R6":     "6",
		"R7":     "7",
		"R8":     "8",
		"R9":     "9",
		"R10":    "10",
		"R11":    "11",
		"R12":    "12",
		"R13":    "13",
		"R14":    "14",
		"R15":    "15",
		"SCREEN": "16384",
		"KBD":    "24576",
		"SP":     "0",
		"LCL":    "1",
		"ARG":    "2",
		"THIS":   "3",
		"THAT":   "4",
	}
	destMap := map[string]string{
		"M":   "001",
		"D":   "010",
		"DM":  "011",
		"MD":  "011",
		"A":   "100",
		"AM":  "101",
		"MA":  "101",
		"AD":  "110",
		"DAS": "110",
		"ADM": "111",
		"AMD": "111",
		"MAD": "111",
		"MDA": "111",
		"DMA": "111",
		"DAM": "111",
	}
	cMap := map[string]string{
		"0":   "0101010",
		"1":   "0111111",
		"-1":  "0111010",
		"D":   "0001100",
		"A":   "0110000",
		"!D":  "0001101",
		"!A":  "0110001",
		"-D":  "0001111",
		"-A":  "0110011",
		"D+1": "0011111",
		"A+1": "0110111",
		"D-1": "0001110",
		"A-1": "0110010",
		"D+A": "0000010",
		"D-A": "0010011",
		"A-D": "0000111",
		"D&A": "0000000",
		"D|A": "0010101",
		"M":   "1110000",
		"!M":  "1110001",
		"-M":  "1110011",
		"M+1": "1110111",
		"M-1": "1110010",
		"D+M": "1000010",
		"D-M": "1010011",
		"M-D": "1000111",
		"D&M": "1000000",
		"D|M": "1010101",
	}
	jMap := map[string]string{
		"JGT": "001",
		"JEQ": "010",
		"JGE": "011",
		"JLT": "100",
		"JNE": "101",
		"JLE": "110",
		"JMP": "111",
	}
	countLines := 0
	startDataRegister := 16
	// Читаем файл построчно в первый раз
	// Прочитываем переменные и добавляем их в таблицу
	for scannerFirst.Scan() {
		line := scannerFirst.Text()
		trimLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimLine, "//") || trimLine == "" {
			continue
		}
		// Обрабатываем А- инструкцию без символов
		if strings.HasPrefix(line, "(") {
			strWithoutPefix, _ := strings.CutPrefix(line, "(")
			swithoutPostFix, _ := strings.CutSuffix(strWithoutPefix, ")")
			symbolMap[swithoutPostFix] = strconv.FormatInt(int64(countLines), 10)
		} else {
			countLines++
		}
	}
	// Читаем файл построчно второй раз
	for scannerSecond.Scan() {
		line := scannerSecond.Text()
		trimLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimLine, "//") || trimLine == "" {
			continue
		}
		// Обрабатываем А- инструкцию
		if strings.HasPrefix(trimLine, "@") {
			strWithoutPefix, _ := strings.CutPrefix(trimLine, "@")
			if isAllDigits(strWithoutPefix) {
				number, _ := strconv.ParseInt(strWithoutPefix, 10, 64)
				binaryString := strconv.FormatInt(number, 2)
				binary16Bit := fmt.Sprintf("%016s", binaryString)
				_, err := writer.WriteString(binary16Bit + "\n")
				fmt.Println(err)
			} else {
				value, ok := symbolMap[strWithoutPefix]
				if ok {
					number, _ := strconv.ParseInt(value, 10, 64)
					binaryString := strconv.FormatInt(number, 2)
					binary16Bit := fmt.Sprintf("%016s", binaryString)
					_, err := writer.WriteString(binary16Bit + "\n")
					fmt.Println(err)
				} else {
					symbolMap[strWithoutPefix] = strconv.FormatInt(int64(startDataRegister), 10)
					startDataRegister++
					number, _ := strconv.ParseInt(symbolMap[strWithoutPefix], 10, 64)
					binaryString := strconv.FormatInt(number, 2)
					binary16Bit := fmt.Sprintf("%016s", binaryString)
					_, err := writer.WriteString(binary16Bit + "\n")
					fmt.Println(err)
				}
			}

		} else {
			// Обрабатываем C- инструкцию
			binaryDest := "000"
			if strings.Contains(trimLine, "=") { // обработка присваивания
				splitLine := strings.Split(trimLine, "=")
				binaryDest := destMap[splitLine[0]]
				cInstruct := cMap[splitLine[1]]
				cBinaryCode := "111" + cInstruct + binaryDest + "000"
				_, err := writer.WriteString(cBinaryCode + "\n")
				fmt.Println(err)
			} else if strings.Contains(trimLine, ";") { // обработка прыжка
				splitJmp := strings.Split(trimLine, ";")
				jmpBinary := jMap[splitJmp[1]]
				cInstruct := cMap[splitJmp[0]]
				cBinaryCode := "111" + cInstruct + binaryDest + jmpBinary
				_, err := writer.WriteString(cBinaryCode + "\n")
				fmt.Println(err)
			}
		}

	}
	if err := scannerFirst.Err(); err != nil {
		log.Fatalf("failed to read file: %s", err)
	}
	if err := scannerSecond.Err(); err != nil {
		log.Fatalf("failed to read file: %s", err)
	}
	// Записываем буфер в файл
	if err := writer.Flush(); err != nil {
		log.Fatalf("failed to flush writer: %s", err)
	}
}
func isAllDigits(s string) bool {
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
