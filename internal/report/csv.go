package report

import (
	"encoding/csv"
	"io"
	"strings"

	"parkvisitor/internal/domain"
)

func CSV(records []domain.VisitorRecord) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"id", "batch_id", "name", "company", "host", "visit_date", "status", "notes"}); err != nil {
		return "", err
	}
	for _, record := range records {
		if err := writer.Write([]string{record.ID, record.BatchID, record.Name, record.Company, record.Host, record.VisitDate, record.Status, record.Notes}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func CSVRowCount(data string) int {
	reader := csv.NewReader(strings.NewReader(data))
	count := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count
		}
		count++
	}
	if count > 0 {
		return count - 1
	}
	return 0
}
