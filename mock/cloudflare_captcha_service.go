package mock

import (
	"context"

	"github.com/weeback/grpc-project-template/internal/entity/cloudflare"
	"github.com/weeback/grpc-project-template/internal/model"
)

var CaptchaService cloudflare.CaptchaService = &captchaServiceMock{}

type captchaServiceMock struct {
}

func (mock *captchaServiceMock) VerifyToken(ctx context.Context, remoteIP string, token string) (*model.TurnstileResult, error) {
	// TODO: Implement the mock logic for verifying the token

	// Simulating a successful verification response
	return &model.TurnstileResult{
		Success:     true,
		ChallengeTS: "2023-10-01T00:00:00Z",
		Hostname:    "example.com",
		Action:      "test_action",
		Cdata:       "",
	}, nil
}
