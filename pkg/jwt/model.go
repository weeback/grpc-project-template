package jwt

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type MapClaims struct {
	jwt.RegisteredClaims
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
	Payload   any    `json:"payload"`
}

func (claims *MapClaims) ApplyContext(ctx context.Context, reqId string) context.Context {
	ctx = context.WithValue(ctx, ApiRequestIdKey, reqId)
	ctx = context.WithValue(ctx, SessionIdKey, claims.SessionId)
	ctx = context.WithValue(ctx, UserIdKey, claims.UserId)
	return ctx
}
