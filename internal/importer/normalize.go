package importer

import (
	"strings"
	"unicode"
)

func CleanText(value string) string {
	fields := strings.FieldsFunc(strings.TrimSpace(value), func(r rune) bool { return unicode.IsSpace(r) || r == '\u3000' })
	return strings.Join(fields, " ")
}

func CleanTags(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ';' || unicode.IsSpace(r) })
	result := []string{}
	seen := map[string]bool{}
	for _, part := range parts {
		clean := strings.ToLower(CleanText(part))
		if clean != "" && !seen[clean] {
			seen[clean] = true
			result = append(result, clean)
		}
	}
	return result
}

func CanonicalCompany(value string) string {
	clean := CleanText(value)
	replacements := map[string]string{"有限公司": "", "有限责任公司": ""}
	for from, to := range replacements {
		clean = strings.TrimSuffix(clean, from)
		if to != "" {
			clean = strings.ReplaceAll(clean, from, to)
		}
	}
	return strings.TrimSpace(clean)
}

func CanonicalHost(value string) string {
	return CleanText(strings.TrimPrefix(strings.TrimSpace(value), "@"))
}

func MergeNotes(existing, incoming string) string {
	left, right := CleanText(existing), CleanText(incoming)
	if left == "" {
		return right
	}
	if right == "" || left == right {
		return left
	}
	return left + "; " + right
}

func RowFingerprint(row Row) string {
	return strings.Join([]string{CleanText(row.Value("name")), CanonicalCompany(row.Value("company")), CanonicalHost(row.Value("host")), CleanText(row.Value("visit_date"))}, "|")
}
