package colly

import (
	"github.com/gocolly/colly"
	"github.com/hill-daniel/glass-scraper"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const https = "https://"

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

	c.OnHTML("div.single-company-result", func(e *colly.HTMLElement) {
		company := glass.Company{}
		nameDiv := e.DOM.Find("div.col-9")
		nameH2 := nameDiv.Find("h2")
		nameH2Anchor := nameH2.Find("a")
		companyName := Purge(nameH2Anchor.Text())
		details, ok := nameH2Anchor.Attr("href")
		if ok {
			company.DetailsURL = baseURL + details
		}
		url := e.DOM.Find("span.url").Text()
		company.URL = https + Purge(url)
		company.Name = Purge(companyName)

		contributionSection := e.DOM.Find("div.ei-contributions-count-wrap")
		reviewAnchor := contributionSection.Find("a.reviews")
		reviews := reviewAnchor.Find("span.num").Text()
		company.NumReviews = int(ParseFloat(reviews))

		ratingText := e.DOM.Find("span.bigRating.strong.margRtSm.h2").Text()
		company.Rating = ParseFloat(ratingText)
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
