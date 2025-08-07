package cloudflare

import (
	"github.com/weeback/grpc-project-template/internal/model"
	"context"
)

type CaptchaService interface {
	VerifyToken(ctx context.Context, remoteIP string, token string) (*model.TurnstileResult, error)
}
