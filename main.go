package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Scraper struct {
	db *gorm.DB
}

func init() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

func main() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Deal{})
	s := &Scraper{db: db}
	s.start()
}

func (u *Scraper) start() {
	dealCollector := colly.NewCollector(
		colly.AllowedDomains("www.mydealz.de", "mydealz.de"),
	)
	dealCollector.Limit(&colly.LimitRule{
		DomainGlob: "*httpbin.*",
		Delay:      5 * time.Second,
	})
	var deals []Deal
	dealCollector.OnHTML("article", func(e *colly.HTMLElement) {
		var deal Deal
		dID, err := strconv.Atoi(strings.Split(e.Attr("id"), "_")[1])
		if err != nil {
			log.Fatalln(err)
		}
		deal.ID = dID
		dClasses := e.Attr("class")
		deal.Expired = strings.Contains(dClasses, "thread--expired")
		deal.Price = e.ChildText("span.thread-price")
		deal.Name = e.ChildText(".thread-title>a")
		deal.Link = e.ChildAttr("a.thread-link", "href")
		deals = append(deals, deal)
	})
	dealCollector.OnScraped(func(r *colly.Response) {
		u.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&deals)
		deals = nil
	})
	dealCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	//dealCollector.Visit("https://www.mydealz.de/new?page=1")
	for i := 1; i < (52999 + 1); i++ {
		dealCollector.Visit(fmt.Sprintf("https://www.mydealz.de/new?page=%d", i))
	}
}
