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
	C_UNKNOWN
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
	line := p.Scanner.Text()
	lineWithoutComments := strings.Split(line, "//")[0]
	p.currentCommand = lineWithoutComments
	p.currentCommandType = p.commandType()
}

func (p *Parser) commandType() Commands {
	arithmeticsCommands := []string{"add", "sub", "neg", "and", "eq", "gt", "lt", "or", "not"}
	lineWithoutComments := strings.Split(p.currentCommand, "//")[0]
	switch {
	case slices.Contains(arithmeticsCommands, lineWithoutComments):
		return C_ARITHMETIC
	case strings.Contains(lineWithoutComments, "push"):
		return C_PUSH
	case strings.Contains(lineWithoutComments, "pop"):
		return C_POP
	case strings.Contains(lineWithoutComments, "label"):
		return C_LABEL
	case strings.HasPrefix(lineWithoutComments, "goto"):
		return C_GOTO
	case strings.Contains(lineWithoutComments, "if"):
		return C_IF
	case strings.Contains(lineWithoutComments, "call"):
		return C_CALL
	case strings.Contains(lineWithoutComments, "function"):
		return C_FUNCTION
	case strings.Contains(lineWithoutComments, "return"):
		return C_RETURN
	default:
		return C_UNKNOWN
	}
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
