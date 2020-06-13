package glass

// Company represents the collected data of a company.
type Company struct {
	Name       string
	Rating     float32
	NumReviews int
	URL        string
	DetailsURL string
}

// Collector collects company information
type Collector interface {
	Collect(filter func(company Company) bool) ([]Company, error)
}

// Exporter exports company data.
type Exporter interface {
	Export([]Company) error
}
