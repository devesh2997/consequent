package middleware

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/devesh2997/consequent/contextx"
	"github.com/gin-gonic/gin"
)

var (
	xRequestIDKey = "X-Request-ID"
)

// generator a function type that returns string.
type generator func() string

var (
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
)

func uuid(len int) string {
	bytes := make([]byte, len)
	random.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)[:len]
}

// RequestInfo is a middleware that injects a RequestID, request body and request url into the context of each request.
func RequestInfo(gen generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithRequestID := injectRequestID(c, gen)
		contextWithRequestURL := injectRequestURL(c.Request, contextWithRequestID)
		contextWithRequestBody := injectRequestBody(c, contextWithRequestURL)
		contextWithRequestHeader := injectRequestHeader(c, contextWithRequestBody)

		c.Request = c.Request.WithContext(contextWithRequestHeader)
		c.Next()
	}
}

func injectRequestID(ctx *gin.Context, gen generator) context.Context {
	xRequestID := GetRequestIDFromHeaders(ctx)

	if xRequestID == "" {
		if gen != nil {
			xRequestID = gen()
		} else {
			xRequestID = uuid(16)
		}
	}

	reqContext := ctx.Request.Context()
	contextWithRequestID := contextx.WithRequestID(reqContext, xRequestID)

	return contextWithRequestID
}

func injectRequestURL(req *http.Request, ctx context.Context) context.Context {
	url := req.URL.String()
	contextWithRequestURL := contextx.WithRequestURL(ctx, url)

	return contextWithRequestURL
}

func injectRequestBody(ginCtx *gin.Context, ctxToInjectIn context.Context) context.Context {
	// Read the content
	var bodyBytes []byte
	if ginCtx.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(ginCtx.Request.Body)
	}

	var body interface{}
	var err error
	if len(bodyBytes) > 0 {
		err = json.Unmarshal(bodyBytes, &body)
	}
	if err != nil {
		body = string(bodyBytes)
	}

	contextWithRequestBody := contextx.WithRequestBody(ctxToInjectIn, body)

	// Restore the io.ReadCloser to its original state
	ginCtx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return contextWithRequestBody
}

func injectRequestHeader(ginCtx *gin.Context, ctxToInjectIn context.Context) context.Context {
	header := ginCtx.Request.Header

	contextWithRequestHeader := contextx.WithRequestHeader(ctxToInjectIn, header)

	return contextWithRequestHeader
}

// GetRequestIDFromHeaders returns 'RequestID' from the headers if present.
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(string(xRequestIDKey))
}
