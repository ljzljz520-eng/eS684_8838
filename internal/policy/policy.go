package policy

import (
	"errors"
	"strings"

	"parkvisitor/internal/domain"
)

type Rules struct {
	AllowedCompanies []string
	RequireHost      bool
	MaxBatchSize     int
}

func DefaultRules() Rules {
	return Rules{RequireHost: true, MaxBatchSize: 500, AllowedCompanies: []string{}}
}

func (r Rules) ValidateInput(input domain.VisitorRecord) error {
	if r.RequireHost && strings.TrimSpace(input.Host) == "" {
		return errors.New("host is required by policy")
	}
	if r.MaxBatchSize <= 0 {
		return errors.New("max batch size policy is invalid")
	}
	if len(r.AllowedCompanies) == 0 {
		return nil
	}
	for _, company := range r.AllowedCompanies {
		if strings.EqualFold(company, input.Company) {
			return nil
		}
	}
	return errors.New("company is not allowed")
}

func (r Rules) AcceptBatch(size int) bool { return size > 0 && size <= r.MaxBatchSize }

func (r Rules) NormalizeCompany(value string) string { return strings.TrimSpace(value) }
