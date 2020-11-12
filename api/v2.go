package api

import (
	"net/url"
	"strings"
)

type v2 struct {
	ServiceURL, ServerURL string
}

func (a *v2) LoginURL(option *LoginOption) string {
	var u strings.Builder
	u.WriteString(a.ServerURL)
	u.WriteString("/login?service=")
	if option == nil {
		u.WriteString(url.QueryEscape(a.ServiceURL))
		return u.String()
	}
	if option.CallbackURL != "" {
		u.WriteString(url.QueryEscape(option.CallbackURL))
	} else {
		u.WriteString(url.QueryEscape(a.ServiceURL))
	}
	if option.Renew {
		u.WriteString("&renew=true")
	}
	if option.Gateway {
		u.WriteString("&gateway=true")
	}
	return u.String()
}

func (a *v2) LogoutURL(option *LogoutOption) string {
	var u strings.Builder
	u.WriteString(a.ServerURL)
	u.WriteString("/logout?service=")
	if option == nil {
		u.WriteString(url.QueryEscape(a.ServiceURL))
	} else {
		u.WriteString(url.QueryEscape(option.CallbackURL))
	}
	return u.String()
}

func (a *v2) ValidateURL(option ValidateOption) string {
	return a.validateURL(&option, "/serviceValidate")
}

func (a *v2) ProxyValidateURL(option ValidateOption) string {
	return a.validateURL(&option, "/proxyValidate")
}

func (a *v2) validateURL(option *ValidateOption, endpoint string) string {
	var u strings.Builder
	u.WriteString(a.ServerURL)
	u.WriteString(endpoint)
	u.WriteString("?service=")
	u.WriteString(url.QueryEscape(a.ServiceURL))
	u.WriteString("&ticket=")
	u.WriteString(option.Ticket)
	if option.Renew {
		u.WriteString("&renew=true")
	}
	if option.PgtURL != "" {
		u.WriteString("&pgtUrl=")
		u.WriteString(url.QueryEscape(option.PgtURL))
	}
	if option.Format != "" {
		u.WriteString("&format=")
		u.WriteString(url.QueryEscape(option.PgtURL))
	}
	return u.String()
}

var _ API = &v2{}

func NewAPIv2(serverURL, serviceURL string) API {
	return &v2{
		ServiceURL: serviceURL,
		ServerURL:  serverURL,
	}
}