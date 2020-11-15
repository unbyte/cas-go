package cas

import (
	"errors"
	"github.com/unbyte/cas-go/api"
	"github.com/unbyte/cas-go/internal/utils"
	"github.com/unbyte/cas-go/parser"
	"net/http"
	"time"
)

type Client interface {
	RedirectToLogin(w http.ResponseWriter, r *http.Request)

	// ValidateTicket validates ticket and if success, return parsed Attributes
	ValidateTicket(ticket string) (parser.Attributes, error)

	// ValidateSession validates session (ticket) and if success, return parsed Attributes
	ValidateSession(sessionID string) (parser.Attributes, error)

	// Validate validates session and if success, save session and write cookie to client
	Validate(w http.ResponseWriter, r *http.Request) (parser.Attributes, error)

	API() api.API

	AttributesStore() AttributesStore
}

type client struct {
	apiInstance     api.API
	client          *http.Client
	attributesStore AttributesStore
	preferredFormat string
	sessionManager  SessionManager
}

var _ Client = &client{}

type Option struct {
	ValidateTimeout time.Duration
	PreferredFormat string
	APIInstance     api.API
	AttributesStore AttributesStore
	SessionManager  SessionManager
}

func New(option Option) Client {
	if option.AttributesStore == nil {
		panic("AttributesStore can't be nil")
	}
	if option.SessionManager == nil {
		panic("SessionManager can't be nil")
	}
	return &client{
		apiInstance:     option.APIInstance,
		client:          &http.Client{Timeout: option.ValidateTimeout},
		preferredFormat: option.PreferredFormat,
		attributesStore: option.AttributesStore,
		sessionManager:  option.SessionManager,
	}
}

func (c *client) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.apiInstance.LoginURL(&api.LoginOption{
		CallbackURL: utils.GetCallbackURLFromRequest(r),
	}), http.StatusFound)
}

func (c *client) ValidateTicket(ticket string) (parser.Attributes, error) {
	u := c.apiInstance.ValidateURL(api.ValidateOption{
		Ticket: ticket,
		Format: c.preferredFormat,
	})

	body, header, err := utils.GetRequest(c.client, u)
	if err != nil {
		return nil, err
	}
	// parser
	ct := header.Get("Content-Type")
	p, ok := c.apiInstance.GetParser(ct)
	if !ok {
		return nil, errors.New("no parser for content-type " + ct)
	}
	// result
	r, success := p(body)
	if !success {
		return r, r.FailureReason()
	}
	return r, nil
}

func (c *client) ValidateSession(sessionID string) (parser.Attributes, error) {
	if a, ok := c.attributesStore.Get(sessionID); ok {
		return a, nil
	}
	return nil, errors.New("it's a new session")
}

func (c *client) Validate(w http.ResponseWriter, r *http.Request) (parser.Attributes, error) {
	if ticket := r.URL.Query().Get("ticket"); ticket != "" {
		a, err := c.ValidateTicket(ticket)
		if err != nil {
			return a, err
		}
		// success, save
		sessionID := utils.GenerateSessionID()
		// ignore errors
		//if err = c.attributesStore.Set(sessionID, a); err != nil {
		//	return a, err
		//}
		//if err := c.sessionManager.Set(w, r, sessionID); err != nil {
		//	return a, err
		//}
		_ = c.attributesStore.Set(sessionID, a)
		_ = c.sessionManager.Set(w, r, sessionID)
		return a, nil
	}
	sessionID, ok := c.sessionManager.Get(r)
	if !ok {
		return nil, errors.New("it's an empty session")
	}
	return c.ValidateSession(sessionID)
}

func (c *client) AttributesStore() AttributesStore {
	return c.attributesStore
}

func (c *client) API() api.API {
	return c.apiInstance
}
