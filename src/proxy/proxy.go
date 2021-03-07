package proxy

import (
	"errors"
	"src/src/framework"
	"time"
)

type Proxy struct {
	url       string
	timestamp int
}

type ProxyCtx struct {
	proxies []Proxy
}

var modules = []framework.Module{rawProxy}

func (ctx *ProxyCtx) ResultCb(result interface{}) {

}

func (ctx *ProxyCtx) IsResultValid() bool {
	return true
}

func (ctx *ProxyCtx) Init(config interface{}) {
	frameworkCtx := framework.NewFramework(ctx)
	ret := frameworkCtx.Run(modules, config)
	proxies = append(proxies, ret.([]Proxy)...)
}

// ProxyWaitTime timeout
const ProxyWaitTime = 5

func (ctx *ProxyCtx) isProxyAvailable(proxy *Proxy) bool {
	if proxy.timestamp == 0 {
		return true
	}
	curTimestamp := time.Now().Second()
	if curTimestamp-proxy.timestamp < ProxyWaitTime {
		return false
	}
	return true
}

var NoAvailableProxy = errors.New("EOF")

func (ctx *ProxyCtx) Get() (string, error) {
	for _, proxy := range ctx.proxies {
		if !ctx.isProxyAvailable(&proxy) {
			continue
		}
		proxy.timestamp = time.Now().Second()
		return proxy.url, nil
	}
	return "", NoAvailableProxy
}
