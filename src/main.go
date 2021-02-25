package main

import (
	"http"
	"log"
)

func logInit() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	logInit()
	log.Println("start")
	resp := http.HttpGet("http://www.baidu.com")
	log.Println(len(resp))
}
