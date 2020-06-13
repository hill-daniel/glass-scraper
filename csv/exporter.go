package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/hill-daniel/glass-scraper"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// Exporter exports company data to a csv file.
type Exporter struct {
}

// NewExporter creates a new exporter.
func NewExporter() *Exporter {
	return &Exporter{}
}

// Export exports company data to a csv file. If successful, a csv file will be created with a timestamp prefix. E.g. 1592064202896750000_companies_result.csv.
// Will have the following columns: Company Name; Rating; Number of Reviews; URL; Details URL
func (*Exporter) Export(companies []glass.Company) error {
	fileName := fmt.Sprintf("%d_companies_result.csv", time.Now().UnixNano())
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error(errors.Wrap(err, "failed to close file"))
		}
	}()
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	err = writeHeader(writer)
	if err != nil {
		return err
	}

	err = writeRows(companies, writer)
	if err != nil {
		return err
	}
	return nil
}

func writeRows(companies []glass.Company, writer *csv.Writer) error {
	for _, c := range companies {
		record := make([]string, 5)
		record[0] = c.Name
		record[1] = fmt.Sprintf("%.1f", c.Rating)
		record[2] = fmt.Sprintf("%d", c.NumReviews)
		record[3] = c.URL
		record[4] = c.DetailsURL
		err := writer.Write(record)

		if err != nil {
			return errors.Wrap(err, "failed to write to file")
		}
	}
	return nil
}

func writeHeader(writer *csv.Writer) error {
	header := make([]string, 5)
	header[0] = "Company Name"
	header[1] = "Rating"
	header[2] = "Number of Reviews"
	header[3] = "URL"
	header[4] = "Details URL"
	err := writer.Write(header)
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	return nil
}
