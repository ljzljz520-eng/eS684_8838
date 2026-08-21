package domain

const (
	StatusImported    = "imported"
	StatusValidated   = "validated"
	StatusNeedsReview = "needs_review"
	StatusApproved    = "approved"
	StatusRejected    = "rejected"
	StatusPublished   = "published"
	StatusArchived    = "archived"
	BatchDraft        = "draft"
	BatchValidated    = "validated"
	BatchConfirmed    = "confirmed"
	BatchPublished    = "published"
	BatchArchived     = "archived"
	TaskOpen          = "open"
	TaskDone          = "done"
)

func ValidRecordStatus(value string) bool {
	switch value {
	case StatusImported, StatusValidated, StatusNeedsReview, StatusApproved, StatusRejected, StatusPublished, StatusArchived:
		return true
	}
	return false
}

func ValidBatchState(value string) bool {
	switch value {
	case BatchDraft, BatchValidated, BatchConfirmed, BatchPublished, BatchArchived:
		return true
	}
	return false
}

func CanReview(state string) bool { return state == BatchValidated || state == BatchConfirmed }

func CanPublish(state string) bool { return state == BatchConfirmed || state == BatchValidated }

func CanArchive(state string) bool { return state == BatchPublished || state == BatchConfirmed }

func TaskStateValid(state string) bool { return state == TaskOpen || state == TaskDone }
