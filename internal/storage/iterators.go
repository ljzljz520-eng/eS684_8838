package storage

import (
	"sort"

	"parkvisitor/internal/domain"
)

func sortVisitors(records []domain.VisitorRecord, by string) {
	sort.SliceStable(records, func(i, j int) bool {
		switch by {
		case "company":
			return records[i].Company < records[j].Company
		case "status":
			return records[i].Status < records[j].Status
		case "updated":
			return records[i].UpdatedAt < records[j].UpdatedAt
		default:
			return records[i].ID < records[j].ID
		}
	})
}

func sortBatches(batches []domain.ImportBatch) {
	sort.SliceStable(batches, func(i, j int) bool {
		if batches[i].BusinessDate == batches[j].BusinessDate {
			return batches[i].ID < batches[j].ID
		}
		return batches[i].BusinessDate < batches[j].BusinessDate
	})
}

func copyVisitors(records []domain.VisitorRecord) []domain.VisitorRecord {
	result := make([]domain.VisitorRecord, len(records))
	copy(result, records)
	for i := range result {
		result[i].Tags = append([]string{}, result[i].Tags...)
	}
	return result
}

func copyBatches(batches []domain.ImportBatch) []domain.ImportBatch {
	result := make([]domain.ImportBatch, len(batches))
	copy(result, batches)
	return result
}

func visitorIDs(records []domain.VisitorRecord) []string {
	result := make([]string, 0, len(records))
	for _, record := range records {
		result = append(result, record.ID)
	}
	sort.Strings(result)
	return result
}
