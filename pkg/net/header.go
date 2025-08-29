package net

const (
	corsAllowOriginHeader      string = "Access-Control-Allow-Origin"
	corsExposeHeadersHeader    string = "Access-Control-Expose-Headers"
	corsMaxAgeHeader           string = "Access-Control-Max-Age"
	corsAllowMethodsHeader     string = "Access-Control-Allow-Methods"
	corsAllowHeadersHeader     string = "Access-Control-Allow-Headers"
	corsAllowCredentialsHeader string = "Access-Control-Allow-Credentials"
	corsRequestMethodHeader    string = "Access-Control-Request-Method"
	corsRequestHeadersHeader   string = "Access-Control-Request-Headers"
	corsOriginHeader           string = "Origin"
	corsVaryHeader             string = "Vary"

	headerOrigin        string = "Origin"
	headerUserAgent     string = "User-Agent"
	headerContentType   string = "Content-Type"
	headerAuthorization string = "Authorization"

	xApiClientId       string = "X-Client-Id"
	xApiRequestId      string = "X-Request-Id"
	xApiMoreError      string = "X-More-Error"
	xApiServiceAccount string = "X-Service-Account"
)
