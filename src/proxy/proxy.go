package proxy

import (
	"src/src/framework"
)

type ProxyInfo struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	Anonymity string `json:"anonymity"`
}

var proxies []ProxyInfo

func init() {
	rawProxy := NewRawProxy()
	modules := []framework.Module{rawProxy}
	rets := framework.Run(modules)
	for _, ret := range rets {
		proxy := ret.(ProxyInfo)
		if proxy.Anonymity == "high_anonymous" {
			proxies = append(proxies, proxy)
		}
	}
}

type RawProxy struct {
	url string
}

func NewRawProxy() *RawProxy {
	proxy := &RawProxy{}
	return proxy
}

func (proxy *RawProxy) Init() {}

func (proxy *RawProxy) NextUrl() (url string, done bool, err error) {
	return proxy.url, true, nil
}

func (proxy *RawProxy) GetHtmlType() string {
	return "json"
}

func (proxy *RawProxy) WithProxy() bool {
	return false
}

func (proxy *RawProxy) IsValid(interface{}) bool {
	return true
}

func (proxy *RawProxy) OnGetHtml(html string) interface{} {
	return html
}
