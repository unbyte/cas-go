package parser

import (
	"strings"
)

type Parser func([]byte) (Attributes Attributes, success bool)

var (
	parsers = map[string]Parser{
		"application/xml": parseXML,
	}
)

func GetParser(contentType string) (Parser, bool) {
	p, ok := parsers[strings.Split(contentType, ";")[0]]
	return p, ok
}

func RegisterParser(contentType string, parser Parser) {
	parsers[contentType] = parser
}
