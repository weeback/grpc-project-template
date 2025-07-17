package jwt

type contextKey string

const (
	// ClaimsKey is the key used to store JWT claims in the context
	ApiRequestIdKey contextKey = "apiRequestId"
	SessionIdKey    contextKey = "sessionIdOfClaims"
	UserIdKey       contextKey = "userIdOfClaims"
)
