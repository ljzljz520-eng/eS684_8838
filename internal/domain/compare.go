package domain

import "strings"

func SameBusinessDay(left, right VisitorRecord) bool { return left.VisitDate == right.VisitDate }

func SameIdentity(left, right VisitorRecord) bool {
	return strings.EqualFold(left.Name, right.Name) && strings.EqualFold(left.Company, right.Company) && strings.EqualFold(left.Host, right.Host)
}

func MergeRecord(base, update VisitorRecord) VisitorRecord {
	result := base
	if update.Name != "" {
		result.Name = update.Name
	}
	if update.Company != "" {
		result.Company = update.Company
	}
	if update.Host != "" {
		result.Host = update.Host
	}
	if update.VisitDate != "" {
		result.VisitDate = update.VisitDate
	}
	if update.Notes != "" {
		result.Notes = update.Notes
	}
	if len(update.Tags) > 0 {
		result.Tags = NormalizeTags(append(result.Tags, update.Tags...))
	}
	if update.UpdatedAt != "" {
		result.UpdatedAt = update.UpdatedAt
	}
	return result
}

func AllowedTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusImported:
		return to == StatusValidated || to == StatusNeedsReview
	case StatusValidated:
		return to == StatusApproved || to == StatusNeedsReview
	case StatusApproved:
		return to == StatusPublished || to == StatusRejected
	case StatusPublished:
		return to == StatusArchived
	case StatusArchived:
		return to == StatusPublished
	}
	return false
}
