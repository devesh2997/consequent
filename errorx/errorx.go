package errorx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"

	"go.mongodb.org/mongo-driver/bson"
)

// The maximum number of stackframes on any error.
var MaxStackDepth = 50

type ErrorCoder interface {
	ErrorCode() int
}
type Stacker interface {
	// Stack returns the callstack formatted the same way that go does
	// in runtime/debug.Stack()
	Stack() string
}

type stacker struct {
	stack  []uintptr
	frames []StackFrame
}

func newStacker() *stacker {
	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(3, stack[:])

	return &stacker{
		stack: stack[:length],
	}
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (err *stacker) Stack() string {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.String()
}

// StackFrames returns an array of frames containing information about the
// stack.
func (err *stacker) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

type QueryError struct {
	*stacker
	Query string
	err   error
	args  []interface{}
}

func NewQueryError(Query string, Err error, args ...interface{}) QueryError {
	return QueryError{
		newStacker(),
		Query,
		Err,
		args,
	}
}

func (err QueryError) Error() string {
	return "error encountered in query"
}

func (err QueryError) Unwrap() error {
	return fmt.Errorf("Error in query : %s |  with args : %v | error = %w", err.Query, err.args, err.err)
}

type APICallError struct {
	*stacker
	URL string
	err error
}

func NewAPICallError(url string, err error) APICallError {
	return APICallError{newStacker(), url, err}
}

func (err APICallError) Error() string {
	return "error encountered while calling api - " + err.err.Error()
}

func (err APICallError) Unwrap() error {
	return err.err
}

type SystemError struct {
	*stacker
	Code int
	err  error
}

func NewSystemError(Code int, err error) SystemError {
	return SystemError{newStacker(), Code, err}
}

func (err SystemError) Error() string {
	return "Some Error Occurred. Please Try Later."
}

func (err SystemError) Unwrap() error {
	return err.err
}

type NotFoundError struct {
	*stacker
	Code             int
	Entity           string
	ResourceLocation string
}

func (notFoundError NotFoundError) ErrorCode() int {
	return notFoundError.Code
}

func NewNotFoundError(code int, entity string, resourceLocation string) NotFoundError {
	return NotFoundError{
		stacker:          newStacker(),
		Code:             code,
		Entity:           entity,
		ResourceLocation: resourceLocation,
	}
}

func (err NotFoundError) Error() string {
	return err.Entity + " not found"
}

func (err NotFoundError) Unwrap() error {
	return fmt.Errorf("%s not found in %s", err.Entity, err.ResourceLocation)
}

type ValidationError struct {
	*stacker
	Code int
	msg  string
}

func (validationError ValidationError) ErrorCode() int {
	return validationError.Code
}

func NewValidationError(Code int, msg string) ValidationError {
	return ValidationError{newStacker(), Code, msg}
}

func (err ValidationError) Error() string {
	return "validation failed : " + err.msg
}

type BusinessError struct {
	*stacker
	Code int
	msg  string
}

func NewBusinessError(Code int, msg string) BusinessError {
	return BusinessError{newStacker(), Code, msg}
}

func (err BusinessError) Error() string {
	return err.msg
}

func (businessError BusinessError) ErrorCode() int {
	return businessError.Code
}

type UnauthorizedError struct {
	*stacker
	Code int
	msg  string
}

func (unauthorizedError UnauthorizedError) ErrorCode() int {
	return unauthorizedError.Code
}

func NewUnauthorizedError(Code int, msg string) UnauthorizedError {
	return UnauthorizedError{newStacker(), Code, msg}
}

func (err UnauthorizedError) Error() string {
	return "unauthorized " + err.msg
}

type UnmarshallingError struct {
	*stacker
	unmarshallerType string
	data             []byte
	err              error
}

func NewBSONUnmarshallingError(err error, data []byte) UnmarshallingError {
	return UnmarshallingError{
		stacker:          newStacker(),
		unmarshallerType: "bson",
		err:              err,
		data:             data,
	}
}

func NewJSONUnmarshallingError(err error, data []byte) UnmarshallingError {
	return UnmarshallingError{
		stacker:          newStacker(),
		unmarshallerType: "json",
		err:              err,
		data:             data,
	}
}

func (err UnmarshallingError) Error() string {
	var e error
	var parsed interface{}
	if err.unmarshallerType == "bson" {
		e = bson.Unmarshal(err.data, &parsed)
	} else if err.unmarshallerType == "json" {
		e = json.Unmarshal(err.data, &parsed)
	}

	if e != nil {
		return fmt.Sprintf("%s unmarshalling error : %v | data : %s", err.unmarshallerType, err.err, string(err.data))
	}

	return fmt.Sprintf("%s unmarshalling error : %v | data : %v", err.unmarshallerType, err.err, parsed)
}

func (err UnmarshallingError) Unwrap() error {
	return err.err
}

// FullError appends the Error() result of err and each of it's wrapped error. Uses ------> as default separator between errors.
func FullError(err error) string {
	return getFullErrorMessageWithSeparator(err, "------>")
}

// FullErrorWithSeparator appends the Error() result of err and each of it's wrapped error. Uses the separtor provided between errors.
func FullErrorWithSeparator(err error, separator string) string {
	return getFullErrorMessageWithSeparator(err, separator)
}

func getFullErrorMessageWithSeparator(err error, sep string) string {
	errMessage := err.Error()

	unwrappedError := errors.Unwrap(err)
	for unwrappedError != nil {
		errMessage = fmt.Sprintf("%s %s %s", errMessage, sep, unwrappedError.Error())
		unwrappedError = errors.Unwrap(unwrappedError)
	}

	return errMessage
}

func GetErrorCode(err error) int {
	if errorCoder, ok := err.(ErrorCoder); ok {
		return errorCoder.ErrorCode()
	}

	return 0
}
