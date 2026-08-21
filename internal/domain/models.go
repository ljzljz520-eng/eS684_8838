package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type VisitorRecord struct {
	ID        string   `json:"id"`
	BatchID   string   `json:"batch_id"`
	Name      string   `json:"name"`
	Company   string   `json:"company"`
	Host      string   `json:"host"`
	VisitDate string   `json:"visit_date"`
	Status    string   `json:"status"`
	Notes     string   `json:"notes"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type ImportBatch struct {
	ID           string `json:"id"`
	Reference    string `json:"reference"`
	Source       string `json:"source"`
	State        string `json:"state"`
	Total        int    `json:"total"`
	Valid        int    `json:"valid"`
	Invalid      int    `json:"invalid"`
	BusinessDate string `json:"business_date"`
	ConfirmedAt  string `json:"confirmed_at"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	BatchID  string `json:"batch_id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
	At       string `json:"at"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Digest   string `json:"digest"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"`
}

type CollaborationTask struct {
	ID       string `json:"id"`
	BatchID  string `json:"batch_id"`
	RecordID string `json:"record_id"`
	Assignee string `json:"assignee"`
	State    string `json:"state"`
	Note     string `json:"note"`
	DueDate  string `json:"due_date"`
}

type Report struct {
	BatchID      string         `json:"batch_id"`
	State        string         `json:"state"`
	Total        int            `json:"total"`
	ByStatus     map[string]int `json:"by_status"`
	ByCompany    map[string]int `json:"by_company"`
	PendingTasks int            `json:"pending_tasks"`
	AuditCount   int            `json:"audit_count"`
}

func NewVisitorRecord(id, batch, name, company, host, date, now string) VisitorRecord {
	return VisitorRecord{ID: id, BatchID: batch, Name: NormalizeName(name), Company: strings.TrimSpace(company), Host: strings.TrimSpace(host), VisitDate: date, Status: StatusImported, CreatedAt: now, UpdatedAt: now, Tags: []string{}}
}

func NormalizeName(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func MakeID(prefix, seed string) string {
	sum := sha256.Sum256([]byte(prefix + ":" + seed))
	return prefix + "-" + hex.EncodeToString(sum[:])[:16]
}

func (r VisitorRecord) IsComplete() bool {
	return r.ID != "" && r.BatchID != "" && r.Name != "" && r.Company != "" && r.VisitDate != ""
}

func (r VisitorRecord) HasTag(tag string) bool {
	for _, item := range r.Tags {
		if item == tag {
			return true
		}
	}
	return false
}

func (r VisitorRecord) AddTag(tag string) VisitorRecord {
	if strings.TrimSpace(tag) == "" || r.HasTag(tag) {
		return r
	}
	r.Tags = append(r.Tags, strings.TrimSpace(tag))
	return r
}
