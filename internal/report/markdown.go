package report

import (
	"fmt"
	"sort"
	"strings"

	"parkvisitor/internal/domain"
)

func Markdown(report domain.Report) string {
	var builder strings.Builder
	builder.WriteString("# Batch ")
	builder.WriteString(report.BatchID)
	builder.WriteString("\n\n")
	builder.WriteString("State: ")
	builder.WriteString(report.State)
	builder.WriteString("\n\n")
	builder.WriteString("| Status | Count |\n| --- | ---: |\n")
	statuses := make([]string, 0, len(report.ByStatus))
	for status := range report.ByStatus {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	for _, status := range statuses {
		builder.WriteString("| ")
		builder.WriteString(status)
		builder.WriteString(" | ")
		builder.WriteString(fmt.Sprint(report.ByStatus[status]))
		builder.WriteString(" |\n")
	}
	builder.WriteString("\nPending tasks: ")
	builder.WriteString(fmt.Sprint(report.PendingTasks))
	builder.WriteString("\n")
	return builder.String()
}

func CompanyTable(report domain.Report) string {
	names := make([]string, 0, len(report.ByCompany))
	for name := range report.ByCompany {
		names = append(names, name)
	}
	sort.Strings(names)
	rows := []string{"company,count"}
	for _, name := range names {
		rows = append(rows, name+","+fmt.Sprint(report.ByCompany[name]))
	}
	return strings.Join(rows, "\n")
}
