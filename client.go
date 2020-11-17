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

	// ValidateTicket validates ticket and if success, return data
	ValidateTicket(ticket string) (interface{}, error)

	// ValidateSession validates session (ticket) and if success, return saved data
	ValidateSession(sessionID string) (interface{}, error)

	// Validate validates session and if success, save session and write cookie to client
	Validate(w http.ResponseWriter, r *http.Request) (interface{}, error)

	API() api.API

	Store() Store
}

type client struct {
	apiInstance     api.API
	client          *http.Client
	store           Store
	preferredFormat string
	sessionManager  SessionManager
	resultHandler   parser.ResultHandler
}

var _ Client = &client{}

type Option struct {
	ValidateTimeout time.Duration
	PreferredFormat string
	APIInstance     api.API
	Store           Store
	SessionManager  SessionManager
	ResultHandler   parser.ResultHandler
}

func New(option Option) Client {
	if option.Store == nil {
		option.Store = DefaultStore(0)
	}
	if option.SessionManager == nil {
		option.SessionManager = DefaultSessionManager("")
	}
	return &client{
		apiInstance:     option.APIInstance,
		client:          &http.Client{Timeout: option.ValidateTimeout},
		preferredFormat: option.PreferredFormat,
		store:           option.Store,
		sessionManager:  option.SessionManager,
		resultHandler:   option.ResultHandler,
	}
}

func (c *client) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.apiInstance.LoginURL(&api.LoginOption{
		CallbackURL: utils.GetCallbackURLFromRequest(r),
	}), http.StatusFound)
}

func (c *client) ValidateTicket(ticket string) (interface{}, error) {
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
		return nil, r.FailureReason()
	}
	if c.resultHandler != nil {
		return c.resultHandler(r), nil
	}
	return r, nil
}

func (c *client) ValidateSession(sessionID string) (interface{}, error) {
	if data, ok := c.store.Get(sessionID); ok {
		return data, nil
	}
	return nil, errors.New("it's a new session")
}

func (c *client) Validate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if ticket := r.URL.Query().Get("ticket"); ticket != "" {
		data, err := c.ValidateTicket(ticket)
		if err != nil {
			return data, err
		}
		// success, save
		sessionID := utils.GenerateSessionID()
		// ignore errors
		//if err = c.store.Set(sessionID, data); err != nil {
		//	return data, err
		//}
		//if err := c.sessionManager.Set(w, r, sessionID); err != nil {
		//	return data, err
		//}
		_ = c.store.Set(sessionID, data)
		_ = c.sessionManager.Set(w, r, sessionID)
		return data, nil
	}
	sessionID, ok := c.sessionManager.Get(r)
	if !ok {
		return nil, errors.New("it's an empty session")
	}
	return c.ValidateSession(sessionID)
}

func (c *client) Store() Store {
	return c.store
}

func (c *client) API() api.API {
	return c.apiInstance
}
