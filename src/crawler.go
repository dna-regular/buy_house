package crawler

import "log"

type House struct {
	url      string
	price    int
	location string
	roomType string // 两室一厅/三室一厅
	date     string
}

type crawler interface {
	Init()
	NextHouse() (house House, done bool, err error)
}

var crawls []crawler

func isHouseValid(house House) bool {
	return false
}

func CrawlRegister(crawl crawler) {
	crawls = append(crawls, crawl)
}

func CrawlRun() []House {
	var houses []House
	for _, craw := range crawls {
		craw.Init()
		for {
			house, done, err := craw.NextHouse()
			if err != nil {
				log.Println(err)
			}
			if done {
				break
			}
			if isHouseValid(house) {
				houses = append(houses, house)
			}
		}
	}
	return houses
}
