package domain

import "strings"

var StatusLabels = map[string]string{StatusImported: "已导入", StatusValidated: "已校验", StatusNeedsReview: "待审核", StatusApproved: "已批准", StatusRejected: "已驳回", StatusPublished: "已发布", StatusArchived: "已归档"}

func StatusLabel(status string) string {
	if label, ok := StatusLabels[status]; ok {
		return label
	}
	return "未知状态"
}

func ParseStatus(label string) string {
	clean := strings.TrimSpace(strings.ToLower(label))
	for status, text := range StatusLabels {
		if clean == status || clean == strings.ToLower(text) {
			return status
		}
	}
	return ""
}

func IsTerminalStatus(status string) bool {
	return status == StatusArchived || status == StatusRejected
}

func IsVisibleStatus(status string) bool { return status != StatusArchived }
