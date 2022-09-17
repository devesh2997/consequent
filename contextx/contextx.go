package contextx

import "context"

type contextKey string

var (
	requestIDKey     contextKey = "request_id"
	requestUserKey   contextKey = "request_user"
	bearerTokenKey   contextKey = "bearer_token"
	requestBodyKey   contextKey = "request_body"
	requestHeaderKey contextKey = "request_header"
	requestURLKey    contextKey = "request_url"
)

type RequestUser struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
}

func (user RequestUser) IsPresent() bool {
	return user.ID != 0
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	contextWithRequestID := context.WithValue(ctx, requestIDKey, requestID)

	return contextWithRequestID
}

// GetRequestID returns request id from the given context if present.
func GetRequestID(ctx context.Context) string {
	v := ctx.Value(requestIDKey)

	if requestID, ok := v.(string); ok {
		return requestID
	}

	return ""
}

func WithRequestUser(ctx context.Context, requestUser RequestUser) context.Context {
	contextWithRequestUser := context.WithValue(ctx, requestUserKey, requestUser)

	return contextWithRequestUser
}

func GetRequestUser(ctx context.Context) RequestUser {
	v := ctx.Value(requestUserKey)

	if requestUser, ok := v.(RequestUser); ok {
		return requestUser
	}

	return RequestUser{}
}

func WithRequestHeader(ctx context.Context, header interface{}) context.Context {
	contextWithRequestHeader := context.WithValue(ctx, requestHeaderKey, header)

	return contextWithRequestHeader
}

func GetRequestHeader(ctx context.Context) interface{} {
	v := ctx.Value(requestHeaderKey)

	return v
}

func WithBearerToken(ctx context.Context, bearerToken string) context.Context {
	contextWithBearerToken := context.WithValue(ctx, bearerTokenKey, bearerToken)

	return contextWithBearerToken
}

// GetBearerToken returns bearer token from the given context if present.
func GetBearerToken(ctx context.Context) string {
	v := ctx.Value(bearerTokenKey)

	if bearerToken, ok := v.(string); ok {
		return bearerToken
	}

	return ""
}

func WithRequestBody(ctx context.Context, body interface{}) context.Context {
	contextWithRequestBody := context.WithValue(ctx, requestBodyKey, body)

	return contextWithRequestBody
}

func GetRequestBody(ctx context.Context) interface{} {
	v := ctx.Value(requestBodyKey)

	return v
}

func WithRequestURL(ctx context.Context, requestURL string) context.Context {
	contextWithRequestURL := context.WithValue(ctx, requestURLKey, requestURL)

	return contextWithRequestURL
}

func GetRequestURL(ctx context.Context) string {
	v := ctx.Value(requestURLKey)

	if requestURL, ok := v.(string); ok {
		return requestURL
	}

	return ""
}
