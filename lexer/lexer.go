package lexer

import (
	"strings"
	"unicode"

	"github.com/qw20012/go-json/token"
)

type Lexer struct {
	input []rune // use 'rune' to handle Unicode
	start int
	end   int
	char  rune
}

func NewLexer(input []byte) *Lexer {
	l := &Lexer{input: []rune(string(input))}
	l.readChar()
	l.addBraceIfNeed()
	return l
}

func (l *Lexer) readChar() {
	if l.end < len(l.input) {
		l.char = l.input[l.end]
	} else {
		l.char = 0
	}
	l.start = l.end
	l.end += 1
}

func (l *Lexer) addBraceIfNeed() {
	start := l.start
	end := l.end

	if l.NewToken().Type == token.STRING && l.NewToken().Type == token.COLON {
		str := append([]byte("{"), string(l.input)...)
		str = append(str, "}"...)
		l.input = []rune(string(str))
	}

	l.start = start
	l.end = end
	l.char = l.input[l.end-1]
}

func (l *Lexer) NewToken() token.Token {
	var tok token.Token
	skipWhitespace(l)
	skipComments(l)

	switch l.char {
	case ':':
		tok = token.NewToken(token.COLON, string(l.char))
	case ',':
		tok = token.NewToken(token.COMMA, string(l.char))
	case '{':
		tok = token.NewToken(token.LBRACE, string(l.char))
	case '}':
		tok = token.NewToken(token.RBRACE, string(l.char))
	case '[':
		tok = token.NewToken(token.LBRACKET, string(l.char))
	case ']':
		tok = token.NewToken(token.RBRACKET, string(l.char))
	default:
		if isInteger(l) {
			tok = token.NewToken(token.INTEGER, string(l.input[l.start:l.end]))
		} else if isBoolean(l) {
			tok = token.NewToken(token.BOOLEAN, string(l.input[l.start:l.end]))
		} else if isString(l) {
			str := string(l.input[l.start:l.end])
			str = strings.Trim(str, `"`)
			tok = token.NewToken(token.STRING, str)
		} else if l.char == rune(0) {
			tok = token.NewToken(token.EOF, "")
		} else {
			tok = token.NewToken(token.INVALID, string(l.char))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) PeakToken() token.Token {
	start := l.start
	end := l.end
	tok := l.NewToken()

	l.start = start
	l.end = end
	l.char = l.input[l.end-1]
	return tok
}

func isInteger(l *Lexer) bool {
	if !unicode.IsDigit(l.char) {
		return false
	}

	endIndex := l.end
	for {
		l.char = l.input[endIndex]
		if !unicode.IsDigit(l.char) {
			break
		}
		endIndex += 1
		if endIndex == len(l.input) {
			break
		}
	}

	if !unicode.IsSpace(l.char) && l.char != ']' && l.char != ',' &&
		l.char != '}' && endIndex < len(l.input) {
		return false
	}

	l.end = endIndex

	return true
}

func isBoolean(l *Lexer) bool {
	return isTrueOrFalse(l, []rune("true")) || isTrueOrFalse(l, []rune("false"))
}

func isTrueOrFalse(l *Lexer, value []rune) bool {
	oldEnd := l.end
	oldChar := l.char

	for i := 0; i < len(value) && l.end < len(l.input) && value[i] == l.char; i++ {
		if i == len(value)-1 {
			return true
		}

		l.char = l.input[l.end]
		l.end += 1
	}

	l.end = oldEnd
	l.char = oldChar
	return false
}

func isString(l *Lexer) bool {
	beginQuote := false
	if l.char == '"' {
		beginQuote = true
	}

	for l.end < len(l.input) {
		l.char = l.input[l.end]

		if !beginQuote {

			if unicode.IsSpace(l.char) || l.char == ':' || l.char == ',' ||
				l.char == '}' {

				return true
			}

			if l.end == len(l.input)-1 {
				l.end += 1
				//fmt.Println("test:" + string(l.input[l.start:l.end]))
				return true
			}
		} else if beginQuote && l.input[l.end] == '"' {
			l.end += 1
			return true
		}

		l.end += 1
	}

	return false
}

func skipWhitespace(l *Lexer) {
	for {
		if !unicode.IsSpace(l.char) {
			return
		}
		l.readChar()
	}
}

func skipComments(l *Lexer) {
	if l.char != '/' || (l.input[l.end] != '/' && l.input[l.end] != '*') {
		return
	}

	if l.input[l.end] == '/' {
		l.readChar()
		for l.char != '\n' {
			l.readChar()
		}
	} else if l.input[l.end] == '*' {
		l.readChar()
		for l.char != '*' || l.input[l.end] != '/' {
			l.readChar()
		}
		l.readChar()
		l.readChar()
	}
	skipWhitespace(l)

	//fmt.Println(string(l.char))
}
