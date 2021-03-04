package proxy

import "src/src/framework"

var modules []framework.Module

func init() {
	proxy := NewRawProxy()
	framework.ModuleRegister(modules, proxy)
}

// GetProxies: get proxy list
func GetProxies() []string {
	rets := framework.Run(modules)
	for ret := range rets {

	}
}

type RawProxy struct {
	url string
}

func NewRawProxy() *RawProxy {
	proxy := &RawProxy{}
	return proxy
}

func (proxy *RawProxy) Init() {

}

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
