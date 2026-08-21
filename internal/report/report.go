package report

import (
	"encoding/json"
	"fmt"
	"sort"

	"parkvisitor/internal/domain"
)

func RenderJSON(value any) ([]byte, error) { return json.MarshalIndent(value, "", "  ") }

func RenderSummary(report domain.Report) string {
	companies := make([]string, 0, len(report.ByCompany))
	for company := range report.ByCompany {
		companies = append(companies, company)
	}
	sort.Strings(companies)
	result := fmt.Sprintf("batch=%s state=%s total=%d pending_tasks=%d audits=%d", report.BatchID, report.State, report.Total, report.PendingTasks, report.AuditCount)
	for _, company := range companies {
		result += fmt.Sprintf(" %s=%d", company, report.ByCompany[company])
	}
	return result
}

func StatusOrder() []string {
	return []string{domain.StatusImported, domain.StatusValidated, domain.StatusNeedsReview, domain.StatusApproved, domain.StatusRejected, domain.StatusPublished, domain.StatusArchived}
}
