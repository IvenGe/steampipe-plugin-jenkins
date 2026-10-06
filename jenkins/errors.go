package jenkins

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Patterns that match gojenkins "not found" style errors.
// The library typically returns status code strings ("404") or messages like
// "Build not found" / "No node found" — never the literal "Not found".
var defaultNotFoundPatterns = []string{
	"404",
	"Build not found",
	"No node found",
	"No label found",
}

// isNotFoundErr reports whether err looks like a missing Jenkins resource.
func isNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, pattern := range defaultNotFoundPatterns {
		if strings.Contains(msg, pattern) {
			return true
		}
	}
	return false
}

// isNotFoundError returns an ErrorPredicate for Jenkins API calls.
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	patterns := notFoundErrors
	if len(patterns) == 0 {
		patterns = defaultNotFoundPatterns
	}
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if err == nil {
			return false
		}
		msg := err.Error()
		for _, pattern := range patterns {
			if strings.Contains(msg, pattern) {
				return true
			}
		}
		return false
	}
}
