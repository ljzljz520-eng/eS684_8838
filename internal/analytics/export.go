package analytics

import (
	"fmt"
	"sort"
	"strings"

	"parkvisitor/internal/domain"
)

func QualityLine(records []domain.VisitorRecord) string {
	quality := Evaluate(records)
	return fmt.Sprintf("total=%d complete=%d incomplete=%d rate=%.2f duplicates=%d companies=%d", quality.Total, quality.Complete, quality.Incomplete, quality.CompletionRate, quality.DuplicateNames, quality.Companies)
}

func DateRange(records []domain.VisitorRecord) (string, string) {
	if len(records) == 0 {
		return "", ""
	}
	dates := make([]string, 0, len(records))
	for _, record := range records {
		dates = append(dates, record.VisitDate)
	}
	sort.Strings(dates)
	return dates[0], dates[len(dates)-1]
}

func JoinCompanies(records []domain.VisitorRecord) string {
	companies := TopCompanies(records, 0)
	return strings.Join(companies, ", ")
}
