package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/devesh2997/consequent/contextx"
	"github.com/devesh2997/consequent/identity/domain/services"
	"github.com/devesh2997/consequent/logger"

	"github.com/gin-gonic/gin"
)

var tokenRequiredMessage = "authorization header is required."
var errTokenRequired = errors.New(tokenRequiredMessage)
var errTokenInvalid = errors.New("invalid token provided")

type Tokens struct {
	BearerToken bearerToken `header:"Authorization"` // jwt token
}

func (tokens Tokens) getJWT() string {
	bearerToken := string(tokens.BearerToken)

	s := strings.SplitAfter(bearerToken, "Bearer ")
	if len(s) > 1 {
		return s[1]
	}

	return ""
}

func (tokens Tokens) validate(tokenService services.TokenService) error {
	isBearerTokenPresent := tokens.BearerToken.isPresent()
	if !isBearerTokenPresent {
		return errTokenRequired
	}
	if _, err := tokenService.Validate(tokens.getJWT()); err != nil {
		return err
	}
	_, err := tokens.BearerToken.getRequestUser()
	if err != nil {
		return err
	}

	return nil
}

func respondWithUnauthenticatedError(c *gin.Context, err error) {
	logger.Log.Error(c.Request.Context(), err)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
}

// Authorisation is...
func Authorisation(tokenService services.TokenService) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		var requestTokens Tokens

		if err := gCtx.ShouldBindHeader(&requestTokens); err != nil {
			respondWithUnauthenticatedError(gCtx, errTokenRequired)
			return
		}

		if err := requestTokens.validate(tokenService); err != nil {
			respondWithUnauthenticatedError(gCtx, err)
			return
		}

		saveTokensAndUserToContext(gCtx, requestTokens)

		gCtx.Next()
	}
}

func saveTokensAndUserToContext(gCtx *gin.Context, tokens Tokens) {
	reqContext := gCtx.Request.Context()

	contextWithBearerToken := contextx.WithBearerToken(reqContext, string(tokens.BearerToken))
	requestUser, _ := tokens.BearerToken.getRequestUser()
	contextWithUser := contextx.WithRequestUser(contextWithBearerToken, *requestUser)

	gCtx.Request = gCtx.Request.WithContext(contextWithUser)
}

type bearerToken string

func (t bearerToken) isPresent() bool {
	return t != ""
}

func (t bearerToken) getRequestUser() (*contextx.RequestUser, error) {
	jwt := strings.Replace(string(t), "Bearer ", "", 1)
	base64EncodedPayload := strings.Split(jwt, ".")[1]

	payloadBytes, err := base64.RawURLEncoding.DecodeString(base64EncodedPayload)
	if err != nil {
		return nil, err
	}

	type userPayload struct {
		ID     int64  `json:"id"`
		Email  string `json:"email"`
		Mobile string `json:"mobile"`
	}
	payload := struct {
		User userPayload `json:"usr"`
	}{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, err
	}

	if payload.User.ID == 0 {
		return nil, errTokenInvalid
	}

	up := payload.User
	return &contextx.RequestUser{
		ID:     up.ID,
		Mobile: up.Mobile,
		Email:  up.Email,
	}, nil
}
