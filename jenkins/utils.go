package jenkins

import "strings"

const (
	titleDescription = "The title of the resource."
)

// splitFullName splits a Jenkins full name like "folder/sub/job" into the
// leaf name and parent path segments expected by gojenkins.GetJob.
func splitFullName(fullName string) (name string, parents []string) {
	parts := strings.Split(fullName, "/")
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			clean = append(clean, part)
		}
	}
	if len(clean) == 0 {
		return "", nil
	}
	return clean[len(clean)-1], clean[:len(clean)-1]
}

// int64FromMap reads a numeric value that may arrive as int64, int, or float64
// (JSON numbers often decode as float64).
func int64FromMap(m map[string]interface{}, key string) (int64, bool) {
	if m == nil {
		return 0, false
	}
	switch v := m[key].(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}
