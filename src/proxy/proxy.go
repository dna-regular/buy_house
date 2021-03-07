package proxy

import (
	"errors"
	"src/src/framework"
	"time"
)

// Proxy proxy
type Proxy struct {
	url       string
	timestamp int
}

// Context proxy context
type Context struct {
	proxies []Proxy
}

var modules = []framework.Module{rawProxy}

// ResultCb cb
func (ctx *Context) ResultCb(result interface{}) {
	ctx.proxies = append(ctx.proxies, result.([]Proxy)...)
}

// IsResultValid check result valid
func (ctx *Context) IsResultValid(result interface{}) bool {
	return true
}

// New new ctx
func New() *Context {
	return &Context{}
}

// Init init
func (ctx *Context) Init(config interface{}) {
	frameworkCtx := framework.NewFramework(ctx, &modules, config)
	frameworkCtx.Run()
}

// ProxyWaitTime timeout
const ProxyWaitTime = 5

func (ctx *Context) isProxyAvailable(proxy *Proxy) bool {
	if proxy.timestamp == 0 {
		return true
	}
	curTimestamp := time.Now().Second()
	if curTimestamp-proxy.timestamp < ProxyWaitTime {
		return false
	}
	return true
}

// ErrNoAvailableProxy err
var ErrNoAvailableProxy = errors.New("EOF")

// Get get
func (ctx *Context) Get() (string, error) {
	for _, proxy := range ctx.proxies {
		if !ctx.isProxyAvailable(&proxy) {
			continue
		}
		proxy.timestamp = time.Now().Second()
		return proxy.url, nil
	}
	return "", ErrNoAvailableProxy
}
