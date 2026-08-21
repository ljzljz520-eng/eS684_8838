package analytics

import (
	"strings"

	"parkvisitor/internal/domain"
)

type Filter struct {
	DateFrom        string
	DateTo          string
	Companies       map[string]bool
	IncludeArchived bool
}

func Match(record domain.VisitorRecord, filter Filter) bool {
	if !filter.IncludeArchived && record.Status == domain.StatusArchived {
		return false
	}
	if filter.DateFrom != "" && record.VisitDate < filter.DateFrom {
		return false
	}
	if filter.DateTo != "" && record.VisitDate > filter.DateTo {
		return false
	}
	if len(filter.Companies) > 0 && !filter.Companies[strings.ToLower(record.Company)] && !filter.Companies[record.Company] {
		return false
	}
	return true
}

func FilterRecords(records []domain.VisitorRecord, filter Filter) []domain.VisitorRecord {
	result := make([]domain.VisitorRecord, 0, len(records))
	for _, record := range records {
		if Match(record, filter) {
			result = append(result, record)
		}
	}
	return result
}

func GroupByDate(records []domain.VisitorRecord) map[string][]domain.VisitorRecord {
	result := map[string][]domain.VisitorRecord{}
	for _, record := range records {
		result[record.VisitDate] = append(result[record.VisitDate], record)
	}
	return result
}

func GroupByHost(records []domain.VisitorRecord) map[string][]domain.VisitorRecord {
	result := map[string][]domain.VisitorRecord{}
	for _, record := range records {
		result[record.Host] = append(result[record.Host], record)
	}
	return result
}
