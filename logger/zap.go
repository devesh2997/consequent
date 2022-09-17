package logger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devesh2997/consequent/errorx"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	// log level
	Level string
	// show caller in log message
	EnableCaller bool
	// error log file paths
	ErrOutputFilePaths []string
	// output log file paths
	OutputFilePaths []string
}

type Field struct {
	Key   string
	Value interface{}
}

// FieldsExtractor extracts fields that need to be logged that are present in the given context
type FieldsExtractor func(ctx context.Context) []Field

type loggerWrapper struct {
	lw              *zap.SugaredLogger
	fieldsExtractor FieldsExtractor
}

func (logger *loggerWrapper) Error(ctx context.Context, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	requestParams = logger.addStackerInfoIfPresent(requestParams, args)

	unwrappedArgs := logger.unwrapArgs(args)
	message := fmt.Sprintf("%v", unwrappedArgs)

	logger.lw.Errorw(message, requestParams...)
}

func (logger *loggerWrapper) unwrapArgs(args []interface{}) []interface{} {
	for i, arg := range args {
		args[i] = logger.unwrapArg(arg)
	}

	return args
}

func (logger *loggerWrapper) unwrapArg(arg interface{}) interface{} {
	if err, ok := arg.(error); ok {
		errMessage := errorx.FullError(err)
		return errMessage
	}

	return arg
}

func (logger *loggerWrapper) Errorf(ctx context.Context, format string, args ...interface{}) {
	unwrappedArgs := logger.unwrapArgs(args)
	message := fmt.Sprintf(format, unwrappedArgs...)
	requestParams := logger.extractFieldsFromContext(ctx)
	requestParams = logger.addStackerInfoIfPresent(requestParams, args)

	logger.lw.Errorw(message, requestParams...)
}

func (logger *loggerWrapper) addStackerInfoIfPresent(requestParams []interface{}, args []interface{}) []interface{} {
	if len(args) == 0 {
		return requestParams
	}
	if stacker, ok := args[0].(errorx.Stacker); ok {
		requestParams = append(requestParams, zap.Reflect("stacker", stacker.Stack()))
	}

	return requestParams
}

func (logger *loggerWrapper) extractFieldsFromContext(ctx context.Context) []interface{} {
	a := []interface{}{}
	if logger.fieldsExtractor == nil {
		return a
	}

	extractedFields := logger.fieldsExtractor(ctx)
	for _, field := range extractedFields {
		a = append(a, zap.Reflect(field.Key, field.Value))
	}

	return a
}

func (logger *loggerWrapper) Fatal(ctx context.Context, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	message := fmt.Sprintf("%v", args)
	logger.lw.Fatalw(message, requestParams...)
}

func (logger *loggerWrapper) Fatalf(ctx context.Context, format string, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	message := fmt.Sprintf(format, args...)
	logger.lw.Fatalw(message, requestParams...)
}

func (logger *loggerWrapper) Info(ctx context.Context, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	message := fmt.Sprintf("%v", args)
	logger.lw.Infow(message, requestParams...)
}

func (logger *loggerWrapper) Infof(ctx context.Context, format string, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	message := fmt.Sprintf(format, args...)
	logger.lw.Infow(message, requestParams...)
}

func (logger *loggerWrapper) Warnf(ctx context.Context, format string, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	message := fmt.Sprintf(format, args...)
	logger.lw.Warnw(message, requestParams...)
}

func (logger *loggerWrapper) Debug(ctx context.Context, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	message := fmt.Sprintf("%v", args)
	logger.lw.Debugw(message, requestParams...)
}

func (logger *loggerWrapper) Debugf(ctx context.Context, format string, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	message := fmt.Sprintf(format, args...)
	logger.lw.Debugw(message, requestParams...)
}

func (logger *loggerWrapper) Println(ctx context.Context, args ...interface{}) {
	requestParams := logger.extractFieldsFromContext(ctx)
	logger.lw.Info(args, requestParams, "\n")
}

func WithFieldsExtractor(lc LogConfig, fieldsExtractor FieldsExtractor) error {
	return registerLog(lc, fieldsExtractor)
}

func RegisterLog(lc LogConfig) error {
	return registerLog(lc, nil)
}

func registerLog(lc LogConfig, fieldsExtractor FieldsExtractor) error {
	zLogger, err := initLog(lc)
	if err != nil {
		return errors.Wrap(err, "RegisterLog")
	}
	defer zLogger.Sync()
	zSugarlog := zLogger.Sugar()
	zSugarlog.Info()

	wrappedLogger := &loggerWrapper{lw: zSugarlog, fieldsExtractor: fieldsExtractor}

	SetLogger(wrappedLogger)
	return nil
}

// initLog create logger
func initLog(lc LogConfig) (zap.Logger, error) {
	rawJSON := []byte(`{
		"level": "info",
		"Development": true,
		"DisableCaller": false,
		"encoding": "json",
		"encoderConfig": {
			"timeKey":        "ts",
			"levelKey":       "level",
			"messageKey":     "msg",
			"nameKey":        "name",
			"stacktraceKey":  "stacktrace",
			"callerKey":      "caller",
			"lineEnding":     "\n",
			"timeEncoder":     "time",
			"levelEncoder":    "lowercaseLevel",
			"durationEncoder": "stringDuration",
			"callerEncoder":   "shortCaller"
		}
	}`)

	var cfg zap.Config
	var zLogger *zap.Logger
	var err error
	//standard configuration
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return *zLogger, errors.Wrap(err, "Unmarshal")
	}
	//customize it from configuration file
	err = customizeLogFromConfig(&cfg, lc)
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}
	zLogger, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}

	zLogger.Debug("logger construction succeeded")

	return *zLogger, nil
}

// customizeLogFromConfig customize log based on parameters from configuration file
func customizeLogFromConfig(cfg *zap.Config, lc LogConfig) error {
	cfg.DisableCaller = !lc.EnableCaller

	cfg.OutputPaths = lc.OutputFilePaths
	cfg.ErrorOutputPaths = lc.ErrOutputFilePaths

	// set log level
	l := zap.NewAtomicLevel().Level()
	err := l.Set(lc.Level)
	if err != nil {
		return errors.Wrap(err, "")
	}
	cfg.Level.SetLevel(l)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return nil
}
