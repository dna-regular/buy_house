package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"src/src/conf"
	"strings"

	"log"
)

type ProxyInfo struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	Anonymity string `json:"anonymity"`
}

type RawProxy struct {
	url string
}

var rawProxy = &RawProxy{}

func (proxy *RawProxy) Init(config interface{}) {
	cnf := config.(*conf.Config)
	proxy.url = cnf.Proxies.RawProxy
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

func (proxy *RawProxy) OnGetHtml(html string) interface{} {
	var proxies []Proxy
	reader := bufio.NewReader(strings.NewReader(html))
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Printf("parse json proxies done!")
			break
		}
		var info ProxyInfo
		err = json.Unmarshal(line, &info)
		if err != nil {
			log.Fatal("unmarsha json error")
			continue
		}
		if info.Type == "https" {
			url := info.Type + "://" + info.Host + ":" + fmt.Sprint(info.Port)
			proxy := Proxy{url: url}
			proxies = append(proxies, proxy)
		}
	}
	//log.Println(proxies)
	return proxies
}
