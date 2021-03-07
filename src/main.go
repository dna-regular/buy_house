package main

import (
	"log"
	"src/src/conf"
	"src/src/proxy"
)

func logInit() {
	log.SetFlags(log.Lshortfile)
}

const confPath = "/usr/local/etc/house_crawler.conf"

func main() {
	logInit()
	log.Println("start")
	config, err := conf.ParseConf(confPath)
	if err != nil {
		log.Println(err)
		return
	}
	proxyCtx := proxy.New()
	proxyCtx.Init(config)
	proxy, err := proxyCtx.Get()
	log.Println(proxy)
}
