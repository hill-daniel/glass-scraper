package colly

import (
	"github.com/gocolly/colly"
	"github.com/hill-daniel/glass-scraper"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Collector collects company data via colly HTML scraping.
type Collector struct {
	c   *colly.Collector
	url string
}

// NewCollector creates a new Collector.
func NewCollector(url string) (*Collector, error) {
	c := colly.NewCollector()
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "glassdoor.*/*",
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create collector")
	}
	return &Collector{c: c, url: url}, nil
}

// Collect collects company data via colly HTML scraping.
func (collector *Collector) Collect(filter func(company glass.Company) bool) ([]glass.Company, error) {
	var lastVisited string
	var companies []glass.Company
	c := collector.c
	baseURL := collector.baseURL()

	c.OnRequest(func(r *colly.Request) {
		url := r.URL.String()
		log.Infof("Visiting %s", url)
	})

	c.OnHTML("div.eiHdrModule", func(e *colly.HTMLElement) {
		company := glass.Company{}
		companyLink := e.DOM.Find("a.tightAll.h2")
		companyName := Purge(companyLink.Text())
		details, ok := companyLink.Attr("href")
		if ok {
			company.DetailsURL = baseURL + details
		}
		url := e.DOM.Find("span.url").Text()
		company.URL = Purge(url)
		company.Name = Purge(companyName)
		ratingText := e.DOM.Find("span.bigRating.strong.margRtSm.h1").Text()
		company.Rating = ParseFloat(ratingText)
		reviewsEl := e.DOM.Find("a.eiCell.cell.reviews")
		reviewsText := reviewsEl.Find("span.num.h2").Text()
		company.NumReviews = int(ParseFloat(reviewsText))

		if filter(company) {
			companies = append(companies, company)
		}
	})

	c.OnHTML("li.next", func(e *colly.HTMLElement) {
		nextHref := e.ChildAttr("a", "href")
		nextPageURL := baseURL + nextHref

		err := c.Visit(nextPageURL)
		if err != nil {
			log.Error(errors.Wrapf(err, "failed to visit page %s, last successful page (if any) was %s. Writing csv so far", nextPageURL, lastVisited))
		}
		lastVisited = nextPageURL
	})

	err := c.Visit(collector.url)
	if err != nil {
		return companies, err
	}
	return companies, nil
}

func (collector *Collector) baseURL() string {
	i := strings.Index(collector.url, "glassdoor")
	j := strings.Index(collector.url[i:len(collector.url)], "/")
	return collector.url[0 : i+j]
}
