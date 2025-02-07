package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Parser struct {
	Scanner            *bufio.Scanner
	currentCommand     string
	currentCommandType Commands
}

type Commands int

const (
	C_ARITHMETIC Commands = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

func NewParser(file io.Reader) *Parser {
	return &Parser{
		Scanner: bufio.NewScanner(file),
	}
}

func (p Parser) hasMoreLines() bool {
	return p.Scanner.Scan()
}

func (p *Parser) advance() {
	p.currentCommand = p.Scanner.Text()
	p.currentCommandType = p.commandType()
}

func (p *Parser) commandType() Commands {
	arithmeticsCommands := []string{"add", "sub", "neg", "and", "eq", "gt", "lt", "or", "not"}
	var result Commands
	if slices.Contains(arithmeticsCommands, p.currentCommand) {
		result = C_ARITHMETIC
	}
	if strings.Contains(p.currentCommand, "push") {
		result = C_PUSH
	}
	if strings.Contains(p.currentCommand, "pop") {
		result = C_POP
	}
	return result
}

func (p *Parser) arg1() (string, error) {
	if p.currentCommandType == C_RETURN {
		return "", fmt.Errorf("не может быть вызвана к команде с типо C_RETURN")
	}
	return strings.Split(p.currentCommand, " ")[0], nil
}
func (p *Parser) arg2() (string, error) {
	if slices.Contains([]Commands{C_PUSH, C_POP, C_FUNCTION, C_CALL}, p.currentCommandType) {
		return strings.Split(p.currentCommand, " ")[1], nil
	}
	return "", fmt.Errorf("не может быть вызвана к команде с данным типом")
}
