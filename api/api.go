package api

type API interface {
	LoginURL(option *LoginOption) string
	LogoutURL(option *LogoutOption) string
	ValidateURL(option ValidateOption) string
	ProxyValidateURL(option ValidateOption) string
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
