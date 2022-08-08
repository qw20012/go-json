package parser

import (
	"fmt"

	"github.com/qw20012/go-json/lexer"
	"github.com/qw20012/go-json/token"
)

type Parser struct {
	Lexer *lexer.Lexer
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{Lexer: l}
}

func (p *Parser) Parse() any {
	tok := p.Lexer.NewToken()
	switch tok.Type {
	case token.STRING:
		return string(tok.Lit)
	case token.INTEGER:
		return string(tok.Lit)
	case token.BOOLEAN:
		return string(tok.Lit)
	case token.LBRACE:
		return parseObject(p)
	case token.LBRACKET:
		return parseArray(p)
	case token.EOF:
		return nil
	}
	return nil
}

func parseArray(p *Parser) any {
	array := []any{}
	tok := p.Lexer.PeakToken()

	if tok.Type == token.RBRACKET {
		return array
	} else {
		array = append(array, p.Parse())
		tok = p.Lexer.NewToken()
		if tok.Type == token.RBRACKET {
			return array
		}
	}

	for {
		array = append(array, p.Parse())
		tok = p.Lexer.NewToken()
		if tok.Type == token.RBRACKET {
			break
		}

		if tok.Type != token.COMMA {
			panic(fmt.Sprintf("was expecting ',' got %s in array parse", string(tok.Lit)))
		}
	}

	return array
}

func parseObject(p *Parser) any {
	object := map[string]any{}
	tok := p.Lexer.NewToken()

	if tok.Type == token.RBRACE { // nothing inside
		return object
	} else {
		key := string(tok.Lit)
		p.Lexer.NewToken() // ':'
		object[key] = p.Parse()
		tok = p.Lexer.NewToken()
		if tok.Type == token.RBRACE {
			return object
		}
	}

	for {
		tok = p.Lexer.NewToken()

		// Allow json last line end with ","
		if tok.Type == token.RBRACE {
			break
		}

		key := string(tok.Lit)

		tok = p.Lexer.NewToken() // ':'
		if tok.Type != token.COLON {
			panic(fmt.Sprintf("was expecting ':' got %s", string(tok.Lit)))
		}

		object[key] = p.Parse()
		tok = p.Lexer.NewToken() // ','

		if tok.Type == token.RBRACE {
			break
		}

		if tok.Type != token.COMMA {
			panic(fmt.Sprintf("was expecting ',' got %s", string(tok.Lit)))
		}
	}

	return object
}
