package storage

var bucketNames = [][]byte{
	[]byte("visitor_records"), []byte("import_batches"), []byte("audit_events"), []byte("attachments"), []byte("collaboration_tasks"),
}

const (
	visitorBucket    = "visitor_records"
	batchBucket      = "import_batches"
	auditBucket      = "audit_events"
	attachmentBucket = "attachments"
	taskBucket       = "collaboration_tasks"
)
