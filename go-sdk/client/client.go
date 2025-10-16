package client

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	ctx         context.Context
	cancel      context.CancelFunc
	url         string
	apiKey      string
	tokenHolder *clientToken
}

func NewSDKClient(par context.Context, url string, apiKey string) (*Client, error) {
	ctx, cancel := context.WithCancel(par)
	c := &Client{
		ctx:    ctx,
		cancel: cancel,
		url:    url,
		apiKey: apiKey,
	}
	resp, err := c.loginByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	c.tokenHolder = &clientToken{
		accessToken: resp.AccessToken,
	}
	go c.refreshToken()
	return c, nil
}

func (c *Client) Close() {
	c.cancel()
}

func (c *Client) refreshToken() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.setupToken()
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) setupToken() {
	for {
		resp, err := c.loginByAPIKey(c.apiKey)
		if err != nil {
			logrus.Errorf("failed to refresh token: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		c.tokenHolder.setToken(resp.AccessToken)
		return
	}
}

type clientToken struct {
	mutex       sync.RWMutex
	accessToken string
}

func (ct *clientToken) getToken() string {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()
	return ct.accessToken
}

func (ct *clientToken) setToken(accessToken string) {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()
	ct.accessToken = accessToken
}
