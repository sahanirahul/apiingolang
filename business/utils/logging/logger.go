package logging

import (
	"apiingolang/activity/middleware/corel"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Fields map[string]interface{}

type ILogger interface {
	WriteLogs(ctx context.Context, msg string, level Level, fields Fields)
	Gin() gin.HandlerFunc
}

type zlogger struct {
	logger *zap.Logger
}

var once sync.Once

// making logger exportable so that it can be used from anywhere for now
var Logger *zlogger

// not required for now, can use this in future to inject logger to other struct as dependency
// in such case the struct will have its own instance of logger
func NewLogger() ILogger {
	once.Do(func() {
		initializeLogger()
	})
	return Logger
}

func (l *zlogger) normalizeFields(fields Fields) {
	for key := range fields {
		if fields[key] == nil {
			delete(fields, key)
			continue
		}
		switch val := fields[key].(type) {
		case fmt.Stringer:
			fields[key] = val.String()
		case error:
			fields[key] = val.Error()
		default:
			b, _ := json.Marshal(val)
			fields[key] = string(b)
		}
	}
}

func (l *zlogger) zapFields(fields Fields) []zapcore.Field {
	if len(fields) == 0 {
		return nil
	}
	var zapFields []zapcore.Field
	for key, val := range fields {
		zapFields = append(zapFields, zap.Any(key, val))
	}
	return zapFields
}

func (l *zlogger) WriteLogs(ctx context.Context, msg string, level Level, fields Fields) {
	// do the logging here
	fields[string(corel.RequestIDKey)] = corel.GetRequestIdFromContext(ctx)
	l.normalizeFields(fields)
	zapFields := l.zapFields(fields)
	switch level {
	case InfoLevel:
		//do info logging
		l.logger.Info(msg, zapFields...)
	case ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case PanicLevel:
		l.logger.Panic(msg, zapFields...)
	}
}

func createFile(filepath string) {
	// Create a file
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func initializeLogger() error {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	today := time.Now().Format("2006-01-02")
	logdir := os.Getenv("LOGDIR")

	if len(logdir) == 0 {
		pwd, _ := os.Getwd()
		logdir := path.Join(pwd, "logs")
		if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
			log.Fatal(err)
		}

	}
	logfilePath := path.Join(path.Join(logdir, today))
	createFile(logfilePath)
	logFile, err := os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		// will log debuglevel in file
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		// will log infolevel on console
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)
	Logger = &zlogger{
		logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)),
	}
	return nil
}

// GinLogger returns a gin.HandlerFunc middleware
func (l *zlogger) Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var fields = Fields{}
		var level Level
		level = InfoLevel
		fields["time"] = start
		defer l.WriteLogs(c, "http_request_init", level, fields)
		c.Next()
		stop := time.Since(start)
		fields["latency"] = int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		code := c.Writer.Status()

		fields["statusCode"] = code
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		fields["dataLength"] = dataLength

		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
			level = ErrorLevel
		} else if code > 499 {
			level = ErrorLevel
		} else if code > 399 {
			level = WarnLevel
		}
	}
}
