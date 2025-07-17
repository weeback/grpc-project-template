package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github/weeback/grpc-project-template/internal/entity/cloudflare"
	"github/weeback/grpc-project-template/internal/model"
)

func NewCaptchaService(secretKey, url string) cloudflare.CaptchaService {
	return &CaptchaService{
		Url:       url,
		SecretKey: secretKey,
	}
}

type CaptchaService struct {
	Url       string
	SecretKey string
}

func (ins *CaptchaService) VerifyToken(ctx context.Context, remoteIP, token string) (*model.TurnstileResult, error) {

	var (
		result model.TurnstileResult

		form = url.Values{
			formKeySecret: {ins.SecretKey},
			formKeyToken:  {token},
		}
		challengesURL = DefaultTurnstileVerifyURL
	)

	if remoteIP != "" {
		form.Set(formKeyRemoteIP, remoteIP)
	}
	if ins.Url != "" && ins.Url != "default" {
		challengesURL = ins.Url
	}
	// Create a new HTTP request with the context
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, challengesURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	// Set the content type to application/x-www-form-urlencoded
	request.Header.Set("Content-Type", DefaultContentType)
	// Set the User-Agent header if needed
	request.Header.Set("User-Agent", "Mozilla/5.0 (compatible; CaptchaService/1.0)")
	// Set the Accept header to application/json
	request.Header.Set("Accept", "application/json")
	// Set the remote IP if provided
	if remoteIP != "" {
		request.Header.Set("X-Forwarded-For", remoteIP)
	}
	// Send the request
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		fmt.Printf("%+v\n", map[string]any{
			"func":       "VerifyToken",
			"status":     resp.Status,
			"statusCode": resp.StatusCode,
			"remoteIP":   remoteIP,
			"token":      token,
			"url":        ins.Url,
			"form":       form,
			"error":      err,
			"response":   result,
		})
		if err := Body.Close(); err != nil {
			// Log the error if needed, but do not return it
			// as we are already handling the response.
		}
	}(resp.Body)

	// Check if the response status code is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Read the response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	// Decode the JSON response into the result struct
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	// Check if the verification was successful
	if result.Success {
		return &result, nil
	}
	// If the verification failed, return an error with the error codes
	return &result, fmt.Errorf("turnstile token verification failed")
}
