package crawler

import (
	"log"
)

type House struct {
	url      string
	price    int
	location string
	roomType string // 两室一厅/三室一厅
	date     string
}

type crawler interface {
	Init()
	NextUrl() (url string, done bool, err error)
	GetHouseInPage(html string) (House, err)
}

var crawls []crawler

func isHouseValid(house House) bool {
	return false
}

func CrawlRegister(crawl crawler) {
	crawls = append(crawls, crawl)
}

func GetHouses() []House {
	var houses []House
	for _, crawl := range crawls {
		crawl.Init()
		for {
			url, done, err := crawl.NextUrl()
			if err != nil {
				log.Println(err)
				continue
			}
			if done {
				break
			}
			house, err := crawl.GetHouseInPage(html)
			if err != nil {
				log.Println(err)
				continue
			}
			if isHouseValid(house) {
				houses = append(houses, house)
			}
		}
	}
	return houses
}
