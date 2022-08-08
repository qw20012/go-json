package lexer

import (
	"strings"
	"testing"

	"github.com/qw20012/go-json/token"
)

func TestLexer(t *testing.T) {
	input := `{
				    "glossary": {
				        "title": "example glossary",
						"GlossDiv": {
				            "title": "S",
							"GlossList": {
				                "GlossEntry": {
									"GlossTerm": "Standard Generalized Markup Language",
									"Abbrev": "ISO 8879:1986",
									"GlossDef": {
				                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
										"GlossSeeAlso": ["GML", "XML"]
				                    },
									"GlossSee": "markup"
				                }
				            },
				            "Nums": 5245243
				        }
				    }
				}`

	tests := []struct {
		typ token.Type
		lit string
	}{
		{token.LBRACE, "{"},
		{token.STRING, "\"glossary\""},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.STRING, "\"title\""},
		{token.COLON, ":"},
		{token.STRING, "\"example glossary\""},
		{token.COMMA, ","},
		{token.STRING, "\"GlossDiv\""},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.STRING, "\"title\""},
		{token.COLON, ":"},
		{token.STRING, "\"S\""},
		{token.COMMA, ","},
		{token.STRING, "\"GlossList\""},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.STRING, "\"GlossEntry\""},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.STRING, "\"GlossTerm\""},
		{token.COLON, ":"},
		{token.STRING, "\"Standard Generalized Markup Language\""},
		{token.COMMA, ","},
		{token.STRING, "\"Abbrev\""},
		{token.COLON, ":"},
		{token.STRING, "\"ISO 8879:1986\""},
		{token.COMMA, ","},
		{token.STRING, "\"GlossDef\""},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.STRING, "\"para\""},
		{token.COLON, ":"},
		{token.STRING, "\"A meta-markup language, used to create markup languages such as DocBook.\""},
		{token.COMMA, ","},
		{token.STRING, "\"GlossSeeAlso\""},
		{token.COLON, ":"},
		{token.LBRACKET, "["},
		{token.STRING, "\"GML\""},
		{token.COMMA, ","},
		{token.STRING, "\"XML\""},
		{token.RBRACKET, "]"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.STRING, "\"GlossSee\""},
		{token.COLON, ":"},
		{token.STRING, "\"markup\""},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.STRING, "\"Nums\""},
		{token.COLON, ":"},
		{token.INTEGER, "5245243"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := NewLexer([]byte(input))

	for i, test := range tests {
		tok := l.NewToken()
		if test.typ != tok.Type {
			t.Fatalf("On test[%d], expected Type=%s, Got=%s", i, test.typ, tok.Type)
		}

		if strings.Trim(test.lit, "\"") != string(tok.Lit) {
			t.Fatalf("On test[%d], expected Literal=%s, Got=%s", i, test.lit, string(tok.Lit))
		}
	}
}

func TestPeakToken(t *testing.T) {
	input := `"name":"value"`
	l := NewLexer([]byte(input))
	if l.NewToken().Type != token.LBRACE {
		t.Fatalf("TestPeakToken Without Begin Brace type failed " + string(l.PeakToken().Type))
	}
	if l.PeakToken().Type != token.STRING {
		t.Fatalf("TestPeakToken type failed " + string(l.PeakToken().Type))
	}
	if string(l.PeakToken().Lit) != `name` {
		t.Fatalf("TestPeakToken lit failed " + string(l.PeakToken().Lit))
	}

	inputWithoutQuote := `name:value`
	lexer := NewLexer([]byte(inputWithoutQuote))
	lexer.NewToken()
	if lexer.PeakToken().Type != token.STRING {
		t.Fatalf("TestPeakToken without quote type failed " + string(lexer.PeakToken().Type))
	}
	if string(lexer.PeakToken().Lit) != `name` {
		t.Fatalf("TestPeakToken without quote lit failed " + string(lexer.PeakToken().Lit))
	}

	inputWithBoolean := `true`
	lexerWithBoolean := NewLexer([]byte(inputWithBoolean))

	if lexerWithBoolean.PeakToken().Type != token.BOOLEAN {
		t.Fatalf("TestPeakToken with boolean type failed ")
	}
	if string(lexerWithBoolean.PeakToken().Lit) != `true` {
		t.Fatalf("TestPeakToken with boolean lit failed " + string(lexerWithBoolean.PeakToken().Lit) + ".")
	}

	inputWithLineComments := `// Comments
true `
	lexerWithLineComments := NewLexer([]byte(inputWithLineComments))
	if lexerWithLineComments.PeakToken().Type != token.BOOLEAN {
		t.Fatalf("TestPeakToken with line comments type failed " + string(lexerWithLineComments.PeakToken().Type))
	}
	if string(lexerWithLineComments.PeakToken().Lit) != `true` {
		t.Fatalf("TestPeakToken with line comments lit failed " + string(lexerWithLineComments.PeakToken().Lit) + ".")
	}

	inputWithBlockComments := `/*
	first line
	second line
	*/ false `
	lexerWithBlockComments := NewLexer([]byte(inputWithBlockComments))
	if lexerWithBlockComments.PeakToken().Type != token.BOOLEAN {
		t.Fatalf("TestPeakToken with block comments type failed " + string(lexerWithBlockComments.PeakToken().Type))
	}
	if string(lexerWithBlockComments.PeakToken().Lit) != `false` {
		t.Fatalf("TestPeakToken with block comments lit failed " + string(lexerWithBlockComments.PeakToken().Lit) + ".")
	}

}
