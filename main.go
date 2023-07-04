package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type TokenKind int

const (
	TokenInvalid TokenKind = iota
	TokenNumber
	TokenIdentifier
	TokenOperator
	TokenPunctuation
	TokenKeyword
)

type Token struct {
	Kind  TokenKind
	Value string
}

type Lexer struct {
	input   string
	pos     int
	tokens  []Token
	current Token
}

func NewLexer(input string) *Lexer {
	lexer := &Lexer{
		input:  input,
		pos:    0,
		tokens: make([]Token, 0),
	}
	return lexer
}

func (l *Lexer) Lex() []Token {
	for l.pos < len(l.input) {
		if unicode.IsSpace(rune(l.input[l.pos])) {
			l.pos++
			continue
		}

		if unicode.IsDigit(rune(l.input[l.pos])) {
			l.readNumber()
		} else if unicode.IsLetter(rune(l.input[l.pos])) {
			l.readIdentifier()
		} else {
			l.readSymbol()
		}
	}

	return l.tokens
}

func (l *Lexer) readNumber() {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '.') {
		l.pos++
	}
	l.emitToken(TokenNumber, l.input[start:l.pos])
}

func (l *Lexer) readIdentifier() {
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '_') {
		l.pos++
	}
	l.emitToken(TokenIdentifier, l.input[start:l.pos])
}

func (l *Lexer) readSymbol() {
	start := l.pos
	switch l.input[l.pos] {
	case '+', '-', '*', '/', '(', ')', '{', '}', '=', ';':
		l.pos++
	default:
		l.pos++
		return
	}
	l.emitToken(TokenPunctuation, l.input[start:l.pos])
}

func (l *Lexer) emitToken(kind TokenKind, value string) {
	token := Token{
		Kind:  kind,
		Value: value,
	}
	l.tokens = append(l.tokens, token)
}

type Parser struct {
	tokens  []Token
	current Token
	pos     int
}

func NewParser(tokens []Token) *Parser {
	parser := &Parser{
		tokens:  tokens,
		current: Token{},
		pos:     0,
	}
	return parser
}

func (p *Parser) Parse() {
	p.advance()
	for p.current.Kind != TokenInvalid {
		p.statement()
	}
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
		p.pos++
	} else {
		p.current = Token{Kind: TokenInvalid}
	}
}

func (p *Parser) statement() {
	if p.current.Kind == TokenIdentifier {
		p.assignmentStatement()
	} else {
		p.error("Unexpected token")
	}
}

func (p *Parser) assignmentStatement() {
	identifier := p.current.Value
	p.advance()
	if p.current.Value != "=" {
		p.error("Expected '='")
		return
	}
	p.advance()
	value := p.expression()
	fmt.Printf("Assign %s = %f\n", identifier, value)
}

func (p *Parser) expression() float64 {
	result := p.term()

	for p.current.Kind == TokenOperator && (p.current.Value == "+" || p.current.Value == "-") {
		operator := p.current.Value
		p.advance()
		term := p.term()
		switch operator {
		case "+":
			result += term
		case "-":
			result -= term
		}
	}

	return result
}

func (p *Parser) term() float64 {
	result := p.factor()

	for p.current.Kind == TokenOperator && (p.current.Value == "*" || p.current.Value == "/") {
		operator := p.current.Value
		p.advance()
		factor := p.factor()
		switch operator {
		case "*":
			result *= factor
		case "/":
			result /= factor
		}
	}

	return result
}

func (p *Parser) factor() float64 {
	var result float64

	if p.current.Kind == TokenNumber {
		result, _ = strconv.ParseFloat(p.current.Value, 64)
		p.advance()
	} else if p.current.Kind == TokenIdentifier {
		// TODO: Handle variables
		p.error("Variable not supported")
	} else {
		p.error("Invalid expression")
	}

	return result
}

func (p *Parser) error(message string) {
	fmt.Printf("Parse error: %s\n", message)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter code: ")
	code, _ := reader.ReadString('\n')

	lexer := NewLexer(code)
	tokens := lexer.Lex()

	parser := NewParser(tokens)
	parser.Parse()
}
