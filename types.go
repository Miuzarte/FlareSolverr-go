package FlareSolverr

import (
	"net/http"
	"strings"
	"time"
)

type Param = string

const (
	PARAM_CMD                 Param = "cmd" // [CMD_REQUEST_GET] | [CMD_REQUEST_POST] ...
	PARAM_URL                 Param = "url"
	PARAM_SESSION             Param = "session"
	PARAM_SESSION_TTL_MINUTES Param = "session_ttl_minutes" // int
	PARAM_MAX_TIMEOUT         Param = "maxTimeout"          // int
	PARAM_COOKIES             Param = "cookies"
	PARAM_RETURN_ONLY_COOKIES Param = "returnOnlyCookies" // bool
	PARAM_RETURN_SCREENSHOT   Param = "returnScreenshot"  // bool
	PARAM_PROXY               Param = "proxy"
	PARAM_WAIT_IN_SECONDS     Param = "waitInSeconds" // int
	PARAM_POST_DATA           Param = "postData"      // string // application/x-www-form-urlencoded
)

type Cmd = string

const (
	CMD_REQUEST_GET      Cmd = "request.get"
	CMD_REQUEST_POST     Cmd = "request.post"
	CMD_SESSIONS_LIST    Cmd = "sessions.list"
	CMD_SESSIONS_CREATE  Cmd = "sessions.create"
	CMD_SESSIONS_DESTROY Cmd = "sessions.destroy"
)

const RESP_STATUS_OK = "ok"

type Client struct {
	Endpoint string
}

type RespBase struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	Version        string `json:"version"`
}

type Response struct {
	RespBase
	Session  string    `json:"session"`
	Sessions []string  `json:"sessions"` // sessions.list
	Solution *Solution `json:"solution"`
}

type Solution struct {
	Url        string            `json:"url"`
	Status     int               `json:"status"`
	Cookies    Cookies           `json:"cookies"`
	UserAgent  string            `json:"userAgent"`
	Headers    map[string]string `json:"headers"`
	Response   string            `json:"response"`
	Screenshot string            `json:"screenshot"` // base64-encoded PNG
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	Expiry   int64  `json:"expiry"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"httpOnly"`
	SameSite string `json:"sameSite"`
}

type Cookies []Cookie

func (c *Cookie) ToHttpCookie() *http.Cookie {
	var sameSite http.SameSite
	switch strings.ToLower(c.SameSite) {
	case "lax":
		sameSite = http.SameSiteLaxMode
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteDefaultMode
	}
	return &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  time.Unix(c.Expiry, 0),
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		SameSite: sameSite,
	}
}

func (cs Cookies) ToHttpCookies() []*http.Cookie {
	httpCookies := make([]*http.Cookie, 0, len(cs))
	for _, c := range cs {
		httpCookies = append(httpCookies, c.ToHttpCookie())
	}
	return httpCookies
}
