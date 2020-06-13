package main

import (
	"flag"
	"github.com/hill-daniel/glass-scraper"
	"github.com/hill-daniel/glass-scraper/colly"
	"github.com/hill-daniel/glass-scraper/csv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	lvl, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = log.InfoLevel
	}
	customFormatter := &log.TextFormatter{}
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
}

func main() {
	glassdoorURL := flag.String("glassdoorURL", "", "URL to glassdoor.de city page")
	minRating := flag.Float64("minRating", 0.0, "minimum rating for a company to be collected")
	minReviews := flag.Int("minReviews", 10, "minimum rating for a company to be collected")
	flag.Parse()
	if len(*glassdoorURL) == 0 {
		log.Error("no URL given, exiting...")
		os.Exit(1)
	}

	collector, err := colly.NewCollector(*glassdoorURL)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	minRatingFilter := func(company glass.Company) bool {
		return company.Rating >= float32(*minRating) && company.NumReviews >= *minReviews
	}
	companies, err := collector.Collect(minRatingFilter)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	exporter := csv.NewExporter()
	err = exporter.Export(companies)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info("Done.")
}
