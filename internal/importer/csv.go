package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Row struct {
	Number int
	Values map[string]string
}

var requiredHeaders = []string{"name", "company", "host", "visit_date"}

func ParseCSV(input string) ([]Row, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.TrimLeadingSpace = true
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read headers: %w", err)
	}
	normalized := normalizeHeaders(header)
	if err := checkHeaders(normalized); err != nil {
		return nil, err
	}
	rows := []Row{}
	line := 1
	for {
		values, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		line++
		if readErr != nil {
			return nil, fmt.Errorf("read row %d: %w", line, readErr)
		}
		if len(values) != len(normalized) {
			return nil, fmt.Errorf("row %d has %d fields, expected %d", line, len(values), len(normalized))
		}
		valuesMap := map[string]string{}
		for i, key := range normalized {
			valuesMap[key] = strings.TrimSpace(values[i])
		}
		rows = append(rows, Row{Number: line, Values: valuesMap})
	}
	return rows, nil
}

func normalizeHeaders(headers []string) []string {
	result := make([]string, len(headers))
	for i, header := range headers {
		result[i] = strings.ToLower(strings.TrimSpace(header))
	}
	return result
}

func checkHeaders(headers []string) error {
	seen := map[string]bool{}
	for _, header := range headers {
		if header == "" {
			return errors.New("empty header")
		}
		if seen[header] {
			return fmt.Errorf("duplicate header %s", header)
		}
		seen[header] = true
	}
	for _, required := range requiredHeaders {
		if !seen[required] {
			return fmt.Errorf("missing header %s", required)
		}
	}
	return nil
}

func (r Row) Value(name string) string { return r.Values[strings.ToLower(strings.TrimSpace(name))] }

func (r Row) Complete() bool {
	for _, header := range requiredHeaders {
		if strings.TrimSpace(r.Value(header)) == "" {
			return false
		}
	}
	return true
}

func (r Row) Optional(name string) (string, bool) {
	value, ok := r.Values[strings.ToLower(strings.TrimSpace(name))]
	return value, ok
}
