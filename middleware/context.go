package middleware

import "context"

type contextKey string

const (
	contextKeyUserID    contextKey = "userID"
	contextKeyUserEmail contextKey = "userEmail"
	contextKeyUserRole  contextKey = "userRole"
)

// ContextUserID returns the authenticated user's ID from the request context.
// Returns 0 if not set.
func ContextUserID(ctx context.Context) int {
	if v, ok := ctx.Value(contextKeyUserID).(int); ok {
		return v
	}
	return 0
}

// ContextUserRole returns the authenticated user's role from the request context.
func ContextUserRole(ctx context.Context) string {
	if v, ok := ctx.Value(contextKeyUserRole).(string); ok {
		return v
	}
	return ""
}

// ContextUserEmail returns the authenticated user's email from the request context.
func ContextUserEmail(ctx context.Context) string {
	if v, ok := ctx.Value(contextKeyUserEmail).(string); ok {
		return v
	}
	return ""
}
