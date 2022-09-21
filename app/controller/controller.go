package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/devesh2997/consequent/errorx"
	"github.com/devesh2997/consequent/logger"
	"github.com/gin-gonic/gin"
)

const (
	errorCodeKey    = "error_code"
	codeKey         = "code"
	codeSuccess     = "success"
	codeFailed      = "failed"
	dataKey         = "data"
	msgKey          = "msg"
	errorMessageKey = "errorMessage"
)

type Controller struct {
}

func (ctrl Controller) Send(gCtx *gin.Context, data interface{}) {
	gCtx.JSON(http.StatusOK, gin.H{codeKey: codeSuccess, dataKey: data})
}

func (ctrl Controller) SendSuccess(gCtx *gin.Context) {
	gCtx.JSON(http.StatusOK, gin.H{codeKey: codeSuccess})
}

func (ctrl Controller) SendWithError(gCtx *gin.Context, err error) {
	httpStatusCode := ctrl.getHTTPStatusCodeFromError(err)
	errCode := errorx.GetErrorCode(err)

	// passing context of request because that's where request id is stored.
	ctrl.logError(gCtx.Request.Context(), err)

	gCtx.JSON(httpStatusCode, gin.H{codeKey: codeFailed, errorCodeKey: errCode, errorMessageKey: err.Error()})
}

func (ctrl Controller) SendBadRequestError(gCtx *gin.Context, err error) {
	ctrl.SendWithHTTPStatusCodeAndError(gCtx, http.StatusBadRequest, err)
}

func (ctrl Controller) SendWithHTTPStatusCodeAndError(gCtx *gin.Context, httpStatusCode int, err error) {
	errCode := errorx.GetErrorCode(err)

	// passing context of request because that's where request id is stored.
	ctrl.logError(gCtx.Request.Context(), err)

	gCtx.JSON(httpStatusCode, gin.H{codeKey: codeFailed, errorCodeKey: errCode, errorMessageKey: err.Error()})
}

func (ctrl Controller) SendDataWithError(gCtx *gin.Context, data interface{}, err error) {
	// passing context of request because that's where request id is stored.
	if err != nil {
		ctrl.logError(gCtx.Request.Context(), err)
		gCtx.JSON(http.StatusOK, gin.H{codeKey: codeSuccess, dataKey: data, errorMessageKey: err.Error()})
	} else {
		gCtx.JSON(http.StatusOK, gin.H{codeKey: codeSuccess, dataKey: data})
	}
}

func (ctrl Controller) BindQueryAndBody(gCtx *gin.Context, dest interface{}) error {
	if err := gCtx.ShouldBindQuery(dest); err != nil {
		return err
	}

	if err := gCtx.ShouldBind(dest); err != nil {
		return err
	}

	return nil
}

func (ctrl Controller) logError(ctx context.Context, err error) {
	var businessError errorx.BusinessError
	var validationError errorx.ValidationError

	shouldLogError := true

	if errors.As(err, &validationError) {
		shouldLogError = false
	} else if errors.As(err, &businessError) {
		shouldLogError = false
	}

	if shouldLogError {
		logger.Log.Error(ctx, err)
	}
}

func (ctrl Controller) getHTTPStatusCodeFromError(err error) int {
	var notFoundError errorx.NotFoundError
	var queryError errorx.QueryError
	var apiCallError errorx.APICallError
	var unauthorizedError errorx.UnauthorizedError
	var systemError errorx.SystemError

	if errors.As(err, &notFoundError) {
		return http.StatusNotFound
	} else if errors.As(err, &queryError) {
		return http.StatusInternalServerError
	} else if errors.As(err, &apiCallError) {
		return http.StatusInternalServerError
	} else if errors.As(err, &unauthorizedError) {
		return http.StatusUnauthorized
	} else if errors.As(err, &systemError) {
		return http.StatusInternalServerError
	}

	return 200
}
