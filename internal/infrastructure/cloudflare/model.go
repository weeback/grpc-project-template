package cloudflare

const (
	TurnstileCaptcha CaptchaType = "turnstile"

	DefaultTurnstileVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	DefaultContentType        = "application/x-www-form-urlencoded"

	formKeySecret   = "secret"
	formKeyToken    = "response"
	formKeyRemoteIP = "remoteip"
)

type CaptchaType string
