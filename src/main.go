package main

import (
	"log"
)

func logInit() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	logInit()
	log.Println("start")
	resp := HttpGet("http://www.baidu.com")
	log.Println(len(resp))
}
