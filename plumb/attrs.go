package plumb

import (
	"errors"
	"net/http"
	"strings"
	"unicode"
)

const (
	tokenizeStateInWhitespace = iota
	tokenizeStateInToken
	tokenizeStateInQuote
)

var (
	NoEquals          = errors.New("no '=' in attribute string")
	UnterminatedQuote = errors.New("unterminated quote")
)

func tokenize(s string) ([]string, error) {
	state := tokenizeStateInWhitespace
	nextQuoteIsLiteral := false
	result := []string{}
	for _, ch := range s {
		switch state {
		case tokenizeStateInWhitespace:
			if ch == '\'' {
				state = tokenizeStateInQuote
			} else if !unicode.IsSpace(ch) {
				state = tokenizeStateInToken
				result = append(result, string(ch))
			}
		case tokenizeStateInToken:
			if ch == '\'' && nextQuoteIsLiteral {
				result[len(result)-1] += "'"
				state = tokenizeStateInQuote
			} else if ch == '\'' && !nextQuoteIsLiteral {
				state = tokenizeStateInQuote
			} else if unicode.IsSpace(ch) {
				state = tokenizeStateInWhitespace
			} else {
				result[len(result)-1] += string(ch)
			}
			nextQuoteIsLiteral = false
		case tokenizeStateInQuote:
			if ch == '\'' {
				state = tokenizeStateInToken
				nextQuoteIsLiteral = true
			} else {
				result[len(result)-1] += string(ch)
			}
		}
	}
	if state == tokenizeStateInQuote {
		return result, UnterminatedQuote
	}
	return result, nil
}

// ParseAttributes splits s into tokens, honoring standard plan9 quoting and
// tokenization rules, then stores each name=value token into an attribute
// map.
func ParseAttributes(s string) (map[string]string, error) {
	result := make(map[string]string)
	tokens, err := tokenize(s)
	if err != nil {
		return result, err
	}
	for _, token := range tokens {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return result, NoEquals
		}
		result[http.CanonicalHeaderKey(parts[0])] = parts[1]
	}
	return result, nil
}
