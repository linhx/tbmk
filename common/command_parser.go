package common

import (
	"regexp"
)

type Token struct {
	IsVariable bool   `json:"isVariable"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Raw        string `json:"raw"`
}

func TokensParser(tmplStr string) []Token {
	re := regexp.MustCompile("{{(`?)\\.([a-z][a-zA-Z0-9]*)(\\|(.+?))?(`?)}}")

	matches := re.FindAllStringSubmatchIndex(tmplStr, -1)

	var tokens []Token
	previousIndex := 0
	for _, match := range matches {
		if (match[3] - match[2]) != (match[11] - match[10]) {
			panic("Missing open/close `")
		}

		if previousIndex < match[0] {
			var token = Token{IsVariable: false, Value: tmplStr[previousIndex:match[0]]}
			tokens = append(tokens, token)
		}
		var rawMatched = tmplStr[match[0]:match[1]]
		if match[3]-match[2] == 1 { // escaped
			var token = Token{IsVariable: false, Value: tmplStr[match[2]:match[3]], Raw: rawMatched}
			tokens = append(tokens, token)
		} else {
			var token = Token{IsVariable: true, Name: tmplStr[match[4]:match[5]], Raw: rawMatched}
			if match[8] != -1 {
				token.Value = tmplStr[match[8]:match[9]]
			}
			tokens = append(tokens, token)
		}

		previousIndex = match[1]
	}
	if previousIndex < len(tmplStr) {
		var token = Token{IsVariable: false, Value: tmplStr[previousIndex:len(tmplStr)], Raw: tmplStr[previousIndex:len(tmplStr)]}
		tokens = append(tokens, token)
	}

	return tokens
}
