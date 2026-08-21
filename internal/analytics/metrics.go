package analytics

import (
	"sort"

	"parkvisitor/internal/domain"
)

type Quality struct {
	Total          int
	Complete       int
	Incomplete     int
	CompletionRate float64
	DuplicateNames int
	Companies      int
}

func Evaluate(records []domain.VisitorRecord) Quality {
	result := Quality{Total: len(records)}
	names := map[string]int{}
	companies := map[string]bool{}
	for _, record := range records {
		if record.IsComplete() {
			result.Complete++
		} else {
			result.Incomplete++
		}
		names[record.Name]++
		companies[record.Company] = true
	}
	for _, count := range names {
		if count > 1 {
			result.DuplicateNames += count - 1
		}
	}
	if result.Total > 0 {
		result.CompletionRate = float64(result.Complete) / float64(result.Total)
	}
	result.Companies = len(companies)
	return result
}

func StatusCounts(records []domain.VisitorRecord) map[string]int {
	result := map[string]int{}
	for _, record := range records {
		result[record.Status]++
	}
	return result
}

func CompanyCounts(records []domain.VisitorRecord) map[string]int {
	result := map[string]int{}
	for _, record := range records {
		result[record.Company]++
	}
	return result
}

func TopCompanies(records []domain.VisitorRecord, limit int) []string {
	counts := CompanyCounts(records)
	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		if counts[names[i]] == counts[names[j]] {
			return names[i] < names[j]
		}
		return counts[names[i]] > counts[names[j]]
	})
	if limit < 0 {
		limit = 0
	}
	if len(names) > limit && limit > 0 {
		names = names[:limit]
	}
	return names
}
