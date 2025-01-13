package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)
type Operation func() void

func Translate(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Создаем новый файл для записи
	outputFile, err := os.Create("output.hack")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outputFile.Close()

	// Создаем сканер
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)
	func test () {
		fmt.Println("ТЕСТ")
	}

	// • add, sub , neg
// • eq , gt , lt
// • and, or , not
	logicalMap := map[string]Operation{
		"add": test,
		"sub": test,
		"neg": test,
		"eq": test,
		"gt": test,
		"lt": test,
		"and": test,
		"or": test,
		"not": not
	}

	for scannerFirst.Scan() {
		line := scannerFirst.Text()
		trimLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimLine, "//") || trimLine == "" {
			continue
		}
		if fn, exists := logicalMap[line]; exists {
			
		}
	}
}
