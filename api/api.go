package api

import (
	"github.com/unbyte/cas-go/parser"
	"strings"
)

type API interface {
	LoginURL(option *LoginOption) string
	LogoutURL(option *LogoutOption) string
	ValidateURL(option ValidateOption) string
	ProxyValidateURL(option ValidateOption) string

	GetParser(contentType string) (parser.Parser, bool)
	SetParser(contentType string, parser parser.Parser)
}

type LoginOption struct {
	// will overwrite service url
	CallbackURL string

	Renew bool

	Gateway bool
}

type LogoutOption struct {
	// will overwrite service url
	CallbackURL string
}

type ValidateOption struct {
	Ticket string
	Renew  bool
	PgtURL string
	Format string
}

type parsers map[string]parser.Parser

func (ps parsers) GetParser(contentType string) (parser.Parser, bool) {
	p, ok := ps[strings.Split(contentType, ";")[0]]
	return p, ok
}

func (ps parsers) SetParser(contentType string, parser parser.Parser) {
	ps[contentType] = parser
}
