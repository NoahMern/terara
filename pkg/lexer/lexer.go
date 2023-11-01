package lexer

import "errors"

const (
	TokenEOF = iota
	TokenError

	TokenOpenParen
	TokenCloseParen
	TokenOpenBracket
	TokenCloseBracket
	TokenOpenBrace
	TokenCloseBrace

	TokenComma
	TokenColon
	TokenSemicolon
	TokenDot
	TokenDoubleColon
	TokenPipe

	TokenPercent
	TokenPlus
	TokenMinus
	TokenAsterisk
	TokenSlash
	TokenDoubleSlash
	TokenDoubleAsterisk

	TokenBang
	TokenAnd
	TokenOr
	TokenEqual
	TokenNotEqual
	TokenLessThan
	TokenLessThanEqual
	TokenGreaterThan
	TokenGreaterThanEqual

	TokenInc
	TokenDec

	TokenNumber
	TokenParam
	TokenIdent
	TokenString
	TokenBool
	TokenNull
	TokenLet
)

var tokenNames = map[int]string{
	TokenEOF:              "EOF",
	TokenError:            "Error",
	TokenOpenParen:        "OpenParen",
	TokenCloseParen:       "CloseParen",
	TokenOpenBracket:      "OpenBracket",
	TokenCloseBracket:     "CloseBracket",
	TokenOpenBrace:        "OpenBrace",
	TokenCloseBrace:       "CloseBrace",
	TokenComma:            "Comma",
	TokenColon:            "Colon",
	TokenSemicolon:        "Semicolon",
	TokenDot:              "Dot",
	TokenDoubleColon:      "DoubleColon",
	TokenPipe:             "Pipe",
	TokenPercent:          "Percent",
	TokenPlus:             "Plus",
	TokenMinus:            "Minus",
	TokenAsterisk:         "Asterisk",
	TokenSlash:            "Slash",
	TokenDoubleSlash:      "DoubleSlash",
	TokenDoubleAsterisk:   "DoubleAsterisk",
	TokenBang:             "Bang",
	TokenAnd:              "And",
	TokenOr:               "Or",
	TokenEqual:            "Equal",
	TokenNotEqual:         "NotEqual",
	TokenLessThan:         "LessThan",
	TokenLessThanEqual:    "LessThanEqual",
	TokenGreaterThan:      "GreaterThan",
	TokenGreaterThanEqual: "GreaterThanEqual",
	TokenInc:              "Inc",
	TokenDec:              "Dec",
	TokenNumber:           "Number",
	TokenParam:            "Param",
	TokenIdent:            "Ident",
	TokenString:           "String",
	TokenBool:             "Bool",
	TokenNull:             "Null",
	TokenLet:              "Let",
}

type Token struct {
	Type  int
	Value string
	Pos   int
}

func (t Token) String() string {
	return tokenNames[t.Type]
}

func NewToken(typ int, value string, pos int) *Token {
	return &Token{typ, value, pos}
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input, 0}
}

func (l *Lexer) NextToken() (*Token, error) {
	l.consumeWhitespace()
	if l.pos >= len(l.input) {
		return l.consumeToken(TokenEOF, ""), nil
	}
	switch l.input[l.pos] {
	case '(':
		return l.consumeToken(TokenOpenParen, "("), nil
	case ')':
		return l.consumeToken(TokenCloseParen, ")"), nil
	case '[':
		return l.consumeToken(TokenOpenBracket, "["), nil
	case ']':
		return l.consumeToken(TokenCloseBracket, "]"), nil
	case '{':
		return l.consumeToken(TokenOpenBrace, "{"), nil
	case '}':
		return l.consumeToken(TokenCloseBrace, "}"), nil
	case ',':
		return l.consumeToken(TokenComma, ","), nil
	case ':':
		if l.peek() == ':' {
			return l.consumeToken(TokenDoubleColon, "::"), nil
		}
		return l.consumeToken(TokenColon, ":"), nil
	case ';':
		return l.consumeToken(TokenSemicolon, ";"), nil
	case '.':
		return l.consumeToken(TokenDot, "."), nil
	case '|':
		if l.peek() == '|' {
			return l.consumeToken(TokenOr, "||"), nil
		}
		return l.consumeToken(TokenPipe, "|"), nil
	case '%':
		return l.consumeToken(TokenPercent, "%"), nil
	case '+':
		if l.peek() == '+' {
			return l.consumeToken(TokenInc, "++"), nil
		}
		return l.consumeToken(TokenPlus, "+"), nil
	case '-':
		if l.peek() == '-' {
			return l.consumeToken(TokenDec, "--"), nil
		}
		return l.consumeToken(TokenMinus, "-"), nil
	case '*':
		if l.peek() == '*' {
			return l.consumeToken(TokenDoubleAsterisk, "**"), nil
		}
		return l.consumeToken(TokenAsterisk, "*"), nil
	case '/':
		if l.peek() == '/' {
			return l.consumeToken(TokenDoubleSlash, "//"), nil
		}
		return l.consumeToken(TokenSlash, "/"), nil
	case '!':
		if l.peek() == '=' {
			return l.consumeToken(TokenNotEqual, "!="), nil
		}
		return l.consumeToken(TokenBang, "!"), nil
	case '=':
		if l.peek() == '=' {
			return l.consumeToken(TokenEqual, "=="), nil
		}
		return l.consumeToken(TokenEqual, "="), nil
	case '<':
		if l.peek() == '=' {
			return l.consumeToken(TokenLessThanEqual, "<="), nil
		}
		return l.consumeToken(TokenLessThan, "<"), nil
	case '>':
		if l.peek() == '=' {
			return l.consumeToken(TokenGreaterThanEqual, ">="), nil
		}
		return l.consumeToken(TokenGreaterThan, ">"), nil
	case '&':
		if l.peek() == '&' {
			return l.consumeToken(TokenAnd, "&&"), nil
		}
		return nil, errors.New("Unexpected character: &")
	case '"':
		return l.consumeString('"')
	case '\'':
		return l.consumeString('\'')
	case '$':
		// params start with $ and are followed by an identifier
		return l.consumeParam()
	default:
		if isDigit(l.input[l.pos]) {
			return l.consumeNumber()
		} else if isIdentStart(l.input[l.pos]) {
			return l.consumeIdent()
		}
	}
	return nil, errors.New("Unexpected character: " + string(l.input[l.pos]))
}

func (l *Lexer) consumeWhitespace() {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.pos++
		default:
			return
		}
	}
}

func (l *Lexer) peek() byte {
	return l.input[l.pos]
}

func (l *Lexer) consume() byte {
	b := l.input[l.pos]
	l.pos++
	return b
}

func (l *Lexer) consumeToken(typ int, value string) *Token {
	t := NewToken(typ, value, l.pos)
	l.pos += len(value)
	return t
}

func (l *Lexer) consumeString(delim byte) (*Token, error) {
	l.pos++
	start := l.pos
	for l.pos < len(l.input) {
		if l.input[l.pos] == delim {
			t := NewToken(TokenString, unescapeString(l.input[start:l.pos]), start)
			l.pos++
			return t, nil
		}
		// handle escape sequences
		if l.input[l.pos] == '\\' {
			l.pos++
		}
		l.pos++
	}
	return nil, errors.New("Unterminated string")
}

// handle escape sequences by replacing them with their actual values
func unescapeString(s string) string {
	str := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			switch s[i] {
			case 'n':
				str += "\n"
			case 'r':
				str += "\r"
			case 't':
				str += "\t"
			case '\\':
				str += "\\"
			case '\'':
				str += "'"
			case '"':
				str += "\""
			}
		} else {
			str += string(s[i])
		}
	}
	return str
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (l *Lexer) consumeNumber() (*Token, error) {
	start := l.pos
	for l.pos < len(l.input) {
		if !isDigit(l.input[l.pos]) {
			break
		}
		l.pos++
	}
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		l.pos++
		for l.pos < len(l.input) {
			if !isDigit(l.input[l.pos]) {
				break
			}
			l.pos++
		}
	}
	return NewToken(TokenNumber, l.input[start:l.pos], start), nil
}

func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isIdent(b byte) bool {
	return isIdentStart(b) || isDigit(b)
}

func (l *Lexer) consumeIdent() (*Token, error) {
	start := l.pos
	for l.pos < len(l.input) {
		if !isIdent(l.input[l.pos]) {
			break
		}
		l.pos++
	}
	// handle null,true,false,let
	value := l.input[start:l.pos]
	switch value {
	case "null":
		return NewToken(TokenNull, value, start), nil
	case "true", "false":
		return NewToken(TokenBool, value, start), nil
	case "let":
		return NewToken(TokenLet, value, start), nil
	}
	return NewToken(TokenIdent, l.input[start:l.pos], start), nil
}

func (l *Lexer) consumeParam() (*Token, error) {
	if l.pos+1 >= len(l.input) || !isIdentStart(l.input[l.pos+1]) {
		return nil, errors.New("Expected identifier after $")
	}
	l.pos++
	start := l.pos
	for l.pos < len(l.input) {
		if !isIdent(l.input[l.pos]) {
			break
		}
		l.pos++
	}
	return NewToken(TokenParam, l.input[start:l.pos], start), nil
}
