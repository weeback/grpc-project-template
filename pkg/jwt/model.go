package jwt

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/golang-jwt/jwt/v5"
)

type MapClaims struct {
	jwt.RegisteredClaims
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
	Payload   any    `json:"payload"`
}

func (claims *MapClaims) ParsePayload(v proto.Message) error {
	if claims.Payload == nil {
		return nil
	}
	b, err := json.Marshal(claims.Payload)
	if err != nil {
		return err
	}
	return protojson.Unmarshal(b, v)
}

func (claims *MapClaims) ApplyContext(ctx context.Context, reqId string) context.Context {
	ctx = context.WithValue(ctx, ApiRequestIdKey, reqId)
	ctx = context.WithValue(ctx, SessionIdKey, claims.SessionId)
	ctx = context.WithValue(ctx, UserIdKey, claims.UserId)
	return ctx
}
