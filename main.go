package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type company struct {
	name       string
	rating     float32
	numRatings int
	url        string
}

func main() {
	var visited []string
	var companies []company

	c := colly.NewCollector()
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "glassdoor.de/*",
		Delay:       1 * time.Second,
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to create collector"))
	}

	c.OnRequest(func(r *colly.Request) {
		url := r.URL.String()
		fmt.Println("Visiting", url)
	})

	c.OnHTML("div.eiHdrModule", func(e *colly.HTMLElement) {
		temp := company{}
		companyLink := e.DOM.Find("a.tightAll.h2")
		companyName := purge(companyLink.Text())
		url, ok := companyLink.Attr("href")
		if ok {
			temp.url = "https://www.glassdoor.de" + url
		}
		temp.name = purge(companyName)
		ratingText := e.DOM.Find("span.bigRating.strong.margRtSm.h1").Text()
		temp.rating = parseFloat(ratingText)
		reviewsEl := e.DOM.Find("a.eiCell.cell.reviews")
		reviewsText := reviewsEl.Find("span.num.h2").Text()
		temp.numRatings = int(parseFloat(reviewsText))
		companies = append(companies, temp)
	})

	c.OnHTML("li.next", func(e *colly.HTMLElement) {
		nextHref := e.ChildAttr("a", "href")
		nextPageUrl := "https://www.glassdoor.de" + nextHref
		err := c.Visit(nextPageUrl)
		if err != nil {
			var lastPage string
			if len(visited) > 0 {
				lastPage = visited[len(visited)-1]
			}
			fmt.Printf("failed to visit page %s, last successful page (if any) was %s. Writing csv so far", nextPageUrl, lastPage)
			_ = writeCsv(companies)
			panic(errors.Wrap(err, "failed to visit page"))
		}
		visited = append(visited, nextPageUrl)
	})

	err = c.Visit("https://www.glassdoor.de/Bewertungen/berlin-bewertungen-SRCH_IL.0,6_IM1020.htm")
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing csv...")
	err = writeCsv(companies)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done.")
}

func purge(s string) string {
	return strings.TrimSpace(s)
}

func parseFloat(s string) float32 {
	multiplier := 1.0
	f := strings.Replace(s, ",", ".", -1)
	if strings.Contains(f, " Tsd") {
		f = strings.TrimSuffix(f, " Tsd")
		multiplier = 1000.0
	}
	f = purge(f)
	v, err := strconv.ParseFloat(f, 2)
	if err != nil {
		return 0.0
	}
	return float32(v * multiplier)
}

func writeCsv(companies []company) error {
	fileName := fmt.Sprintf("%d_companies_result.csv", time.Now().UnixNano())
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	header := make([]string, 4)
	header[0] = "Company Name"
	header[1] = "Rating"
	header[2] = "Number of Ratings"
	header[3] = "Url"
	err = writer.Write(header)
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	for _, c := range companies {
		record := make([]string, 4)
		record[0] = c.name
		record[1] = fmt.Sprintf("%.1f", c.rating)
		record[2] = fmt.Sprintf("%d", c.numRatings)
		record[3] = c.url
		err := writer.Write(record)

		if err != nil {
			return errors.Wrap(err, "failed to write to file")
		}
	}
	return nil
}
