package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"parkvisitor/internal/clock"
	"parkvisitor/internal/domain"
	"parkvisitor/internal/storage"
)

type App struct {
	store    *storage.Store
	clock    clock.Clock
	mu       sync.Mutex
	sequence int
}

type VisitorInput struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Company   string   `json:"company"`
	Host      string   `json:"host"`
	VisitDate string   `json:"visit_date"`
	Notes     string   `json:"notes"`
	Tags      []string `json:"tags"`
}

type ImportSummary struct {
	Batch   domain.ImportBatch     `json:"batch"`
	Records []domain.VisitorRecord `json:"records"`
	Issues  map[string][]string    `json:"issues"`
}

func NewApp(store *storage.Store, fixed clock.Clock) (*App, error) {
	if store == nil {
		return nil, errors.New("store is required")
	}
	return &App{store: store, clock: fixed}, nil
}

func (a *App) Store() *storage.Store { return a.store }

func (a *App) nextID(prefix string) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sequence++
	return domain.MakeID(prefix, fmt.Sprintf("%s-%d", a.clock.NowString(), a.sequence))
}

func (a *App) audit(batchID, recordID, action, actor, detail string) error {
	event := domain.AuditEvent{ID: a.nextID("audit"), BatchID: batchID, RecordID: recordID, Action: action, Actor: actor, Detail: detail, At: a.clock.NowString()}
	return a.store.SaveAudit(event)
}

func (a *App) parseBusinessDate(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", errors.New("visit date is required")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		parsed, err = time.Parse("2006-01-02", value)
	}
	if err != nil {
		return "", fmt.Errorf("invalid visit date: %w", err)
	}
	return parsed.UTC().Format("2006-01-02"), nil
}

func (a *App) CreateBatch(id, reference, source string) (domain.ImportBatch, error) {
	if strings.TrimSpace(id) == "" {
		return domain.ImportBatch{}, errors.New("batch id is required")
	}
	batch := domain.ImportBatch{ID: id, Reference: reference, Source: source, State: domain.BatchDraft}
	if err := a.store.SaveBatch(batch); err != nil {
		return batch, err
	}
	_ = a.audit(id, "", "batch_created", "system", source)
	return batch, nil
}

func (a *App) ensureBatch(id string) (domain.ImportBatch, error) {
	batch, err := a.store.GetBatch(id)
	if err == nil {
		return batch, nil
	}
	if strings.TrimSpace(id) == "" {
		return batch, err
	}
	return a.CreateBatch(id, id, "api")
}
