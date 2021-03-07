package framework

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Instance instance
type Instance interface {
	IsResultValid() bool
	ResultCb(result interface{})
}

// Framework framework
type Framework struct {
	ctx     Instance
	modules *[]Module
	config  interface{}
}

// Module module
type Module interface {
	Init(config interface{})
	NextUrl() (url string, done bool, err error)
	GetHtmlType() string
	WithProxy() bool
}

// JSONModule module
type JSONModule interface {
	Module
	OnGetHtml(html string) interface{}
}

// HTMLModule html module
type HTMLModule interface {
	Module
	OnPageParsed(doc *goquery.Document) (interface{}, error)
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

var wg sync.WaitGroup

// NewFramework new framework
func NewFramework(ctx Instance, modules *[]Module, config interface{}) *Framework {
	return &Framework{
		ctx:     ctx,
		config:  config,
		modules: modules,
	}
}

func (framework *Framework) fetchHTML(module Module, url string, results chan Result) {
	result := Result{
		module: module,
		url:    url,
	}
	cli := &http.Client{}
	body := []byte{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	defer wg.Done()
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

func (framework *Framework) handleResult(result Result) {
	if result.statusCode != 200 {
		return
	}
	module := result.module.(Module)
	htmlType := module.GetHtmlType()
	var ret interface{}
	if htmlType == "json" {
		ret = result.module.(JSONModule).OnGetHtml(string(result.body))
	}
	if htmlType == "html" {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(result.body))
		if err != nil {
			log.Fatal(err)
			return
		}
		ret, err = result.module.(HTMLModule).OnPageParsed(doc)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if module.IsValid(ret) {
		rets = append(rets, ret)
	}
}

func (framework *Framework) handleHTML(results chan Result, quit chan int) (rets []interface{}) {
	for {
		select {
		case result := <-results:
			framework.handleResult(result)
		case <-quit:
			log.Println("handleHTML quit")
			return
		}
	}

	return rets
}

const maxCh = 100

// Run run
func (framework *Framework) Run(modules []Module, config interface{}) interface{} {
	var res []interface{}
	results := make(chan Result, maxCh)
	quit := make(chan int)
	go framework.handleHTML(results, quit)
	for _, module := range modules {
		module.Init(config)
		for {
			url, done, err := module.NextUrl()
			if err != nil {
				log.Println(err)
				continue
			}
			if done {
				log.Println("fire tasks done")
				break
			}
			wg.Add(1)
			go framework.fetchHTML(module, url, results)
		}
	}
	wg.Wait()
	quit <- 0
	return res
}
