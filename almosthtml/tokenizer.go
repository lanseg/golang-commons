package almosthtml

import (
	"fmt"
	"strings"
	"unicode"
)

// State is a convenient name for the tokenizer state
type State int

const (
	errorState                  = State(0)
	eofState                    = State(1)
	dataState                   = State(2)
	tagState                    = State(3)
	tagContentState             = State(4)
	tagContentSingleQuotedState = State(5)
	tagContentDoubleQuotedState = State(6)

	// Scripts
	tagScriptState = State(7)

	// Comments
	tagCommentState          = State(8)
	tagCommentMaybeState     = State(9)
	tagCommentMaybeDashState = State(10)
)

var (
	names = []string{
		"errorState",
		"eofState",
		"dataState",
		"tagState",
		"tagContentState",
		"tagContentSingleQuotedState",
		"tagContentDoubleQuotedState",
		"tagScriptState",

		// Comments
		"tagCommentState",
		"tagCommentMaybeState",
		"tagCommentMaybeDashState",
	}
)

func (s State) String() string {
	return names[s]
}

type token struct {
	name      string
	values    []string
	tokenType State
	start     int
	end       int
	empty     bool
}

func (t *token) push(s string) {
	t.empty = false
	if t.tokenType != dataState && t.tokenType != tagScriptState && t.name == "" {
		t.name = s
		return
	}
	t.values = append(t.values, s)
}

func (t *token) String() string {
	if t.tokenType == dataState || t.tokenType == tagCommentState {
		return fmt.Sprintf("%-6s [%d:%d] %q", t.tokenType, t.start, t.end, strings.Join(t.values, ", "))
	}
	return fmt.Sprintf("%-6s [%d:%d] %q %q", t.tokenType, t.start, t.end, t.name, strings.Join(t.values, ", "))
}

// Tokenizer
type tokenizer struct {
	state  State
	buffer strings.Builder
	tokens []*token
	pos    int
}

func (t *tokenizer) lastToken() *token {
	l := len(t.tokens)
	if l == 0 {
		return nil
	}
	return t.tokens[len(t.tokens)-1]
}

func (t *tokenizer) setState(s State) {
	value := t.buffer.String()
	if value != "" {
		t.lastToken().push(value)
	}
	t.state = s
	t.buffer.Reset()
}

func (t *tokenizer) feed(r rune) {
	t.buffer.WriteRune(r)
}

func (t *tokenizer) newToken() {
	tok := &token{
		values:    []string{},
		tokenType: t.state,
		start:     t.pos,
		empty:     true,
	}
	lastIndex := len(t.tokens) - 1
	if lastIndex < 0 {
		t.tokens = append(t.tokens, tok)
	} else if t.tokens[lastIndex].empty {
		t.tokens[lastIndex] = tok
	} else {
		t.tokens[lastIndex].end = t.pos
		t.tokens = append(t.tokens, tok)
	}
}

func tokenize(doc string) []*token {
	t := &tokenizer{
		state:  dataState,
		buffer: strings.Builder{},
		pos:    0,
		tokens: []*token{},
	}
	t.newToken()
	asRunes := []rune(doc)
	for i, r := range asRunes {
		t.pos = i
		switch t.state {
		case errorState:
			fmt.Println("errorState")
			t.setState(dataState)
		case dataState:
			if r == '<' {
				t.setState(tagState)
				t.pos = i + 1
				t.newToken()
			} else {
				t.feed(r)
			}
		case tagState:
			if r == '!' {
				t.state = tagCommentMaybeState
				t.feed(r)
			} else if r == '"' {
				t.setState(tagContentDoubleQuotedState)
			} else if r == '\'' {
				t.setState(tagContentSingleQuotedState)
			} else if r == '=' {
				t.setState(tagContentState)
				t.feed(r)
				t.setState(tagState)
			} else if r == '<' {
				t.setState(errorState)
			} else if r == '>' {
				lastToken := t.lastToken()
				if lastToken.empty {
					t.setState(errorState)
				} else {
					t.setState(dataState)
				}
				t.newToken()
			} else if !unicode.IsSpace(r) {
				t.setState(tagContentState)
				t.feed(r)
			}
		case tagContentDoubleQuotedState:
			if r == '"' {
				t.setState(tagState)
			} else {
				t.feed(r)
			}
		case tagContentSingleQuotedState:
			if r == '\'' {
				t.setState(tagState)
			} else {
				t.feed(r)
			}
		case tagContentState:
			if unicode.IsSpace(r) {
				t.setState(tagState)
			} else if r == '>' {
				t.setState(dataState)
				if t.lastToken().name == "script" {
					t.state = tagScriptState
				}
				t.newToken()
			} else if r == '=' {
				t.setState(tagContentState)
				t.feed(r)
				t.setState(tagState)
			} else {
				t.feed(r)
			}
		case tagScriptState:
			// TODO: Rewrite this dirty hack with suffix
			if r == '>' && strings.HasSuffix(t.buffer.String(), "script") {
				t.setState(dataState)
				t.newToken()
			} else {
				t.feed(r)
			}
		case tagCommentMaybeState:
			t.feed(r)
			if r == '-' {
				t.state = tagCommentMaybeDashState
			} else {
				t.state = tagState
			}
		case tagCommentMaybeDashState:
			t.feed(r)
			if r == '-' {
				t.lastToken().tokenType = tagCommentState
				t.setState(tagCommentState)
			} else {
				t.setState(tagState)
			}
		case tagCommentState:
			if r == '>' && strings.HasSuffix(t.buffer.String(), "--") {
				t.setState(dataState)
				t.newToken()
			} else {
				t.feed(r)
			}
		}
	}
	t.setState(eofState)

	lastToken := t.lastToken()
	if lastToken == nil {
		return []*token{}
	}
	if lastToken.empty {
		t.tokens = t.tokens[:len(t.tokens)-1]
	} else {
		lastToken.end = len(asRunes)
	}
	return t.tokens
}
