package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Condo struct {
	name     string
	address  string
	district string
	// 资产期限
	tenure        string
	numberOfUnits int
	// 开发商
	developer string
	facility  facility
}

type facility struct {
	// 膝上泳池
	lapPool bool
	// 浅水池
	wadingPool bool
	// 按摩
	jacuzzi bool
	// 网球场
	tennisCourt bool
	// 读书角
	readingCorner bool
	// 屋顶花园
	rooftopGarden bool
	// 健身区
	fitnessArea bool
	// 俱乐部
	clubHouse bool
	// 健身房
	gymnasium bool
	// 烧烤设备
	bbqPit bool
	// 秘密花园
	secretGarden bool
	// 慢跑道
	joggingTrack bool
	// 蒸汽室
	steamRoom bool
	// 停车场
	carPark bool
	// 安保
	security bool
}

func main() {
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("singaporeexpats.com",
			"www.singaporeexpats.com",
			"condo.singaporeexpats.com"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		//colly.CacheDir("./coursera_cache"),
	)

	detailCollector := c.Clone()

	condos := make([]*Condo, 0)

	c.OnRequest(func(request *colly.Request) {
		log.Println("visiting", request.URL)
	})

	c.OnHTML("a.title_link", func(e *colly.HTMLElement) {
		condoUrl := e.Request.AbsoluteURL(e.Attr("href"))
		if isCondoUrl(condoUrl) {
			detailCollector.Visit(condoUrl)
		}
	})

	detailCollector.OnHTML("div.propcol1", func(element *colly.HTMLElement) {
		log.Println("condo url", element.Request.URL)
		address := element.ChildText(":nth-child(2)")
		log.Println(address)
		condos = append(condos, &Condo{address: address})
	})

	c.Visit("https://condo.singaporeexpats.com/name/0-9")
}

func validateUrl(url string) bool {
	if strings.HasPrefix(url, "https://condo.singaporeexpats.com/condo") ||
		strings.HasPrefix(url, "https://condo.singaporeexpats.com/name") ||
		strings.HasPrefix(url, "https://condo.singaporeexpats.com/1/name") ||
		strings.HasPrefix(url, "https://condo.singaporeexpats.com/2/name") ||
		strings.HasPrefix(url, "https://condo.singaporeexpats.com/3/name") {
		return true
	}

	return false
}

func isCondoUrl(url string) bool {
	if strings.HasPrefix(url, "https://condo.singaporeexpats.com/condo") {
		return true
	}

	return false
}
