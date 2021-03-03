package framework

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Module interface {
	Init()
	NextUrl() (url string, done bool, err error)
	GetHtmlType() string
	withProxy() bool
}

type JsonModule interface {
	Module
	OnGetHtml(html string) interface{}
}

type HtmlModule interface {
	OnPageParsed() interface{}
}

func ModuleRegister(modules []interface{}, module interface{}) {
	modules = append(modules, module)
}

func fetchHtml(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		continue
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		continue
	}
}

func handeHtml() {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		continue
	}
	house, err := module.GetHouseInPage(html)
	if err != nil {
		log.Println(err)
		continue
	}
	if isHouseValid(house) {
		houses = append(houses, house)
	}
}

func Run(modules []Module) []interface{} {
	res := interface{}
	go handleHtml()
	for _, module := range modules {
		htmlType := module.GetHtmlType()
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
			go fetchHtml(url)
		}
	}
	return houses
}
