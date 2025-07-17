package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	issue   = "Bankaool, S.A., Institución de Banca Múltiple"
	subject = "Bankaool, S.A., Institución de Banca Múltiple"
)

func SignWithClaims(key interface{}, payload any, opts ...*Option) (string, error) {

	opt := NewOption()
	for i, op := range opts {
		if op.liveTime > 0 && (i == 0 || op.liveTime < opt.liveTime) {
			opt = opt.SetLiveTime(op.liveTime)
		}
		if i == 0 && op.sessionId != "" {
			opt = opt.SetSessionId(op.sessionId)
		}
		if i == 0 && op.userId != "" {
			opt = opt.SetUserId(op.userId)
		}
	}

	claims := MapClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    issue,
			Subject:   subject,
			Audience:  []string{"LOGIN"},
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(opt.LiveTime())),
		},
		SessionId: opt.SessionId(),
		UserId:    opt.UserId(),
		Payload:   payload,
	}

	var method jwt.SigningMethod
	switch key.(type) {
	case ecdsa.PrivateKey:
		method = jwt.SigningMethodES256
	case ed25519.PrivateKey:
		method = jwt.SigningMethodEdDSA
	default:
		return "", errors.New("only support ecdsa or ed25519 key")
	}

	// Create a new JWT value
	return jwt.NewWithClaims(method, &claims).SignedString(key)
}

func ParseClaims(pub ed25519.PublicKey, str string) (*MapClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(str, &MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); ok {
			return pub, nil
		}
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	})

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		fmt.Println("❌ Token hết hạn")
		return nil, err
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		fmt.Println("❌ Token chưa đến thời điểm hợp lệ (nbf)")
		return nil, err
	case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
		fmt.Println("❌ Token được sử dụng trước khi được phát hành (iat)")
		return nil, err
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		fmt.Println("❌ Chữ ký không hợp lệ")
		return nil, err
	default:
		if err != nil {
			return nil, err
		}
		if claim, ok := parsedToken.Claims.(*MapClaims); ok {
			return claim, nil
		}
		return nil, fmt.Errorf("invalid token claims")
	}
}
