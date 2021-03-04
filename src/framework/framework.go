package framework

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Module module
type Module interface {
	Init()
	NextUrl() (url string, done bool, err error)
	GetHtmlType() string
	WithProxy() bool
	IsValid(interface{}) bool
}

// Proxy proxy
type Proxy interface {
	GetProxy()
}

// JSONModule module
type JSONModule interface {
	Module
	OnGetHtml(html string) interface{}
}

// HTMLModule html module
type HTMLModule interface {
	Module
	OnPageParsed() (interface{}, error)
}

// Result res
type Result struct {
	module     interface{}
	url        string
	statusCode int
	status     string
	body       []byte
}

// ModuleRegister register
func ModuleRegister(modules []Module, module Module) {
	modules = append(modules, module)
}

const userAget = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"

func fetchHTML(module Module, url string, results chan Result) {
	result := Result{
		module: module,
		url:    url,
	}
	cli := &http.Client{}
	body := []byte{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
		return
	}
	if module.WithProxy() {

	}
	req.Header.Add("User-Agent", userAget)
	resp, err := cli.Do(req)
	defer resp.Body.Close()
	result.status = resp.Status
	result.statusCode = resp.StatusCode
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
		return
	}
	if result.body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Fatalf("read body err: %v", err)
		return
	}
	results <- result
}

func handleHTML(results chan Result) (rets []interface{}) {
	for result := range results {
		if result.statusCode != 200 {
			continue
		}
		module := result.module.(Module)
		htmlType := module.GetHtmlType()
		var ret interface{}
		if htmlType == "json" {
			ret = result.module.(JSONModule).OnGetHtml(string(result.body))
		}
		if htmlType == "html" {
			doc, err := goquery.NewDocumentFromReader(result.body)
			if err != nil {
				log.Fatal(err)
				continue
			}
			ret, err = result.module.(HTMLModule).OnPageParsed(doc)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		if module.IsValid(ret) {
			rets = append(rets, ret)
		}
	}
	return rets
}

const maxCh = 100

// Run run
func Run(modules []Module) []interface{} {
	var res []interface{}
	results := make(chan Result, maxCh)
	go handleHTML(results)
	for _, module := range modules {
		module.Init()
		for {
			url, done, err := module.NextUrl()
			if err != nil {
				log.Println(err)
				continue
			}
			if done {
				break
			}
			go fetchHTML(module, url, results)
		}
	}
	return res
}
