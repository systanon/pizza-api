package apperror

const (
	// Common
	CodeValidation     = "VALIDATION_ERROR"
	CodeNotFound       = "NOT_FOUND"
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeForbidden      = "FORBIDDEN"
	CodeInternal       = "INTERNAL_ERROR"
	CodeRateLimit      = "RATE_LIMIT"
	CodeDuplicateEntry = "DUPLICATE_ENTRY"
	CodeConflict       = "CONFLICT"
	// Auth
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeNoRefreshToken     = "NO_REFRESH_TOKEN"
	CodeInvalidRefresh     = "INVALID_REFRESH"
	CodeEmailNotConfirmed  = "EMAIL_NOT_CONFIRMED"
	CodeResetTokenUsed     = "RESET_TOKEN_ALREADY_USED"
	CodeResetTokenExpired  = "RESET_TOKEN_EXPIRED"
	CodeResetTokenNotFound = "RESET_TOKEN_NOT_FOUND"
	// API
	CodeRouteNotFound = "ROUTE_NOT_FOUND"
)
