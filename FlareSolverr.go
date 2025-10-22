package FlareSolverr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
)

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) Get(ctx context.Context, url string, params map[string]any) (*Solution, error) {
	return c.RequestGet(ctx, url, params)
}

func (c *Client) Post(ctx context.Context, url string, postData string, params map[string]any) (*Solution, error) {
	return c.RequestPost(ctx, url, postData, params)
}

// params:
//
//	(*)url
//	session
//	session_ttl_minutes
//	maxTimeout: 60000
//	cookies: [{"name": "cookie1", "value": "value1"}, {"name": "cookie2", "value": "value2"}]
//	returnOnlyCookies: false
//	returnScreenshot: false
//	proxy: {"url": "http://127.0.0.1:7890", "username": "testuser", "password": "testpass"}
//	waitInSeconds: 0 // Useful to allow it to load dynamic content.
func (c *Client) RequestGet(ctx context.Context, url string, params map[string]any) (*Solution, error) {
	var p map[string]any
	if params == nil {
		p = map[string]any{
			PARAM_URL: url,
		}
	} else {
		p = params
		p[PARAM_URL] = url
	}
	resp, err := c.Submit(ctx, CMD_REQUEST_GET, p)
	if err != nil {
		return nil, err
	}
	return resp.Solution, nil
}

// params:
//
//	(*)url
//	postData: "a=b&c=d" // application/x-www-form-urlencoded
//	// other params same as [Client.RequestGet]
func (c *Client) RequestPost(ctx context.Context, url string, postData string, params map[string]any) (*Solution, error) {
	var p map[string]any
	if params == nil {
		p = map[string]any{
			PARAM_URL:       url,
			PARAM_POST_DATA: postData,
		}
	} else {
		p = params
		p[PARAM_URL] = url
		p[PARAM_POST_DATA] = postData
	}
	resp, err := c.Submit(ctx, CMD_REQUEST_POST, p)
	if err != nil {
		return nil, err
	}
	return resp.Solution, nil
}

// params:
//
//	(*)session
//	proxy: {"url": "http://127.0.0.1:7890", "username": "testuser", "password": "testpass"}
func (c *Client) SessionsCreate(ctx context.Context, session string, params map[string]any) error {
	var p map[string]any
	if params == nil {
		p = map[string]any{
			PARAM_SESSION: session,
		}
	} else {
		p = params
		p[PARAM_SESSION] = session
	}
	_, err := c.Submit(ctx, CMD_SESSIONS_CREATE, p)
	if err != nil {
		return err
	}
	return nil
}

// params:
//
//	(*)session
func (c *Client) SessionsDestroy(ctx context.Context, session string) error {
	_, err := c.Submit(ctx, CMD_SESSIONS_DESTROY, map[string]any{
		PARAM_SESSION: session,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SessionsList(ctx context.Context) ([]string, error) {
	resp, err := c.Submit(ctx, CMD_SESSIONS_LIST, nil)
	if err != nil {
		return nil, err
	}
	return resp.Sessions, nil
}

// Submit 直接提交命令和参数至 FlareSolverr, 返回 Response 结构体
func (c *Client) Submit(ctx context.Context, cmd string, params map[string]any) (*Response, error) {
	var body []byte
	if params != nil {
		p := maps.Clone(params) // 避免修改传入的参数
		p[PARAM_CMD] = cmd
		var err error
		body, err = json.Marshal(p)
		if err != nil {
			return nil, err
		}
	} else {
		body = []byte(`{"` + PARAM_CMD + `":"` + cmd + `"}`)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	fsResp := &Response{}
	err = json.NewDecoder(resp.Body).Decode(fsResp)
	if err != nil {
		return nil, err
	}
	if fsResp.Status != RESP_STATUS_OK {
		return fsResp, errors.New(fsResp.Message)
	}
	return fsResp, nil
}
