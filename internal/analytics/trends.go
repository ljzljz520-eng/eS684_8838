package analytics

import (
	"sort"

	"parkvisitor/internal/domain"
)

type DayCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func DailyCounts(records []domain.VisitorRecord) []DayCount {
	grouped := GroupByDate(records)
	dates := make([]string, 0, len(grouped))
	for date := range grouped {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	result := make([]DayCount, 0, len(dates))
	for _, date := range dates {
		result = append(result, DayCount{Date: date, Count: len(grouped[date])})
	}
	return result
}

func PeakDay(records []domain.VisitorRecord) (DayCount, bool) {
	counts := DailyCounts(records)
	if len(counts) == 0 {
		return DayCount{}, false
	}
	peak := counts[0]
	for _, item := range counts[1:] {
		if item.Count > peak.Count || item.Count == peak.Count && item.Date < peak.Date {
			peak = item
		}
	}
	return peak, true
}

func StatusRate(records []domain.VisitorRecord, status string) float64 {
	if len(records) == 0 {
		return 0
	}
	count := 0
	for _, record := range records {
		if record.Status == status {
			count++
		}
	}
	return float64(count) / float64(len(records))
}

func HasQualityIssues(records []domain.VisitorRecord) bool {
	quality := Evaluate(records)
	return quality.Incomplete > 0 || quality.DuplicateNames > 0
}

func SortedStatuses(records []domain.VisitorRecord) []string {
	counts := StatusCounts(records)
	result := make([]string, 0, len(counts))
	for status := range counts {
		result = append(result, status)
	}
	sort.Strings(result)
	return result
}
