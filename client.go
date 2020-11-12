package cas

import (
	"errors"
	"fmt"
	"github.com/unbyte/cas-go/api"
	"github.com/unbyte/cas-go/internal/utils"
	"github.com/unbyte/cas-go/parser"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	RedirectToLogin(w http.ResponseWriter, r *http.Request)

	// ValidateTicket validates ticket and if success, return parsed Attributes
	ValidateTicket(ticket string) (parser.Attributes, error)

	// ValidateSession validates session (ticket) and if success, return parsed Attributes
	ValidateSession(r *http.Cookie) (parser.Attributes, error)

	// Validate validates session and if success, save session and write cookie to client
	Validate(w http.ResponseWriter, r *http.Request, cookieName string) (parser.Attributes, error)

	API() api.API
	SessionStore() SessionStore
	TicketStore() TicketStore
}

type client struct {
	apiInstance     api.API
	client          *http.Client
	sessionStore    SessionStore
	ticketStore     TicketStore
	preferredFormat string
}

var _ Client = &client{}

type Option struct {
	ValidateTimeout time.Duration
	PreferredFormat string
	APIInstance     api.API
	SessionStore    SessionStore
	TicketStore     TicketStore
}

func New(option Option) Client {
	if option.SessionStore == nil {
		panic("sessionStore can't be nil")
	}
	return &client{
		apiInstance:     option.APIInstance,
		client:          &http.Client{Timeout: option.ValidateTimeout},
		preferredFormat: option.PreferredFormat,
		sessionStore:    option.SessionStore,
		ticketStore:     option.TicketStore,
	}
}

func (c *client) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.apiInstance.LoginURL(&api.LoginOption{
		CallbackURL: utils.GetCallbackURLFromRequest(r),
	}), http.StatusFound)
}

func (c *client) ValidateTicket(ticket string) (parser.Attributes, error) {
	if c.ticketStore != nil {
		if a, ok := c.ticketStore.Get(ticket); ok {
			return a, nil
		}
	}

	u := c.apiInstance.ValidateURL(api.ValidateOption{
		Ticket: ticket,
		Format: c.preferredFormat,
	})

	fmt.Println(u)

	body, header, err := utils.GetRequest(c.client, u)
	if err != nil {
		return nil, err
	}
	// parser
	ct := header.Get("Content-Type")
	p, ok := parser.GetParser(ct)
	if !ok {
		return nil, errors.New("no parser for content-type " + ct)
	}
	// result
	r, success := p(body)
	if !success {
		return nil, errors.New("fail")
	}
	if c.ticketStore != nil {
		_ = c.ticketStore.Set(ticket, r)
	}
	return r, nil
}

func (c *client) ValidateSession(r *http.Cookie) (parser.Attributes, error) {
	if t, ok := c.sessionStore.Get(r.Value); ok {
		a, err := c.ValidateTicket(t)
		if err != nil {
			// validate fail
			_ = c.sessionStore.Del(r.Value)
		}
		return a, err
	}
	return nil, errors.New("it's a new session")
}

func (c *client) Validate(w http.ResponseWriter, r *http.Request, cookieName string) (parser.Attributes, error) {
	if ticket := r.URL.Query().Get("ticket"); ticket != "" {
		a, err := c.ValidateTicket(ticket)
		if err != nil {
			return nil, err
		}
		// success, save
		sessionID := utils.GenerateSessionID()
		if err = c.sessionStore.Set(sessionID, ticket); err != nil {
			return nil, err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    url.QueryEscape(sessionID),
			Path:     "/",
			HttpOnly: true,
		})
		return a, nil
	}
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}
	return c.ValidateSession(cookie)
}

func (c *client) SessionStore() SessionStore {
	return c.sessionStore
}

func (c *client) TicketStore() TicketStore {
	return c.ticketStore
}

func (c *client) API() api.API {
	return c.apiInstance
}
