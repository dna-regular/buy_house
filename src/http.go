package main

import (
	"bytes"
	"log"
	"net/http"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"

func HttpGet(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("create req err: ", err)
		return ""
	}
	req.Header.Add("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("get error: ", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("get error: ", url, resp.StatusCode)
		return ""
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024*1024*2))
	buf.ReadFrom(resp.Body)
	return string(buf.Bytes())
}
