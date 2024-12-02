package common

import (
	"regexp"
)

type Token struct {
	Id         int    `json:"id"`
	IsVariable bool   `json:"isVariable"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Raw        string `json:"raw"`
}

func TokensParser(tmplStr string) []Token {
	re := regexp.MustCompile("(\\\\)?({{([a-z][a-zA-Z0-9]*)(\\|(.+?))?}})")

	matches := re.FindAllStringSubmatchIndex(tmplStr, -1)

	var tokens []Token
	previousIndex := 0
	id := 0
	for _, match := range matches {
		if previousIndex < match[0] { // add literal string
			id++
			var token = Token{Id: id, IsVariable: false, Value: tmplStr[previousIndex:match[0]]}
			tokens = append(tokens, token)
		}
		var rawMatched = tmplStr[match[0]:match[1]]
		if match[3]-match[2] == 1 { // escaped
			id++
			var token = Token{Id: id, IsVariable: false, Value: tmplStr[match[4]:match[5]], Raw: rawMatched}
			tokens = append(tokens, token)
		} else {
			id++
			var token = Token{Id: id, IsVariable: true, Name: tmplStr[match[6]:match[7]], Raw: rawMatched}
			if match[10] != -1 {
				token.Value = tmplStr[match[10]:match[11]]
			}
			tokens = append(tokens, token)
		}

		previousIndex = match[1]
	}
	if previousIndex < len(tmplStr) {
		id++
		var token = Token{Id: id, IsVariable: false, Value: tmplStr[previousIndex:], Raw: tmplStr[previousIndex:]}
		tokens = append(tokens, token)
	}

	return tokens
}
