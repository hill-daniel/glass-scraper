# HTML Scraper for Glassdoor
Scrapes companies and exports them to csv. Contains the following information:
* Name
* URL
* Rating
* Number of reviews
* Detailed information URL

## Usage
* `go build cmd/glass-scraper/main.go`
* execute binary with args:
  * `-glassdoorURL` - mandatory; string; The URL to start scraping. Supports currently location pages, e.g. https://www.glassdoor.co.uk/Reviews/london-reviews-SRCH_IL.0,6_IM1035.htm
  * `-minRating` - optional; float; The minimum rating for a company to be collected, e.g. 3.5. Default 0.0.
  * `-minReviews` - optional; int; The minimum number of reviews for a company to be collected, e.g. 100. Default 10.

## Example execution
`./main -glassdoorURL "https://www.glassdoor.co.uk/Reviews/london-reviews-SRCH_IL.0,6_IM1035.htm" -minRating 3.5 -minReviews 50`
 
 ## CSV Export
 Creates a csv file with the following content:
 
 | Company Name  | Rating        | Number of Reviews | URL                        | Details URL   |
 |:------------- |:-------------:|------------------:|:--------------------------- |--------------:|
 | SAP           | 4.7           | 15000             | https://www.sap.com        | ...           |
 | Hubspot       | 4.6           | 1200              | https://www.hubspot.com    | ...           |
 | Productsup    | 4.5           | 53                | https://www.productsup.com | ...           |