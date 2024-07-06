package middleware

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/o-ga09/langchain-go/pkg/logger"
)

type RequestId string

func AddID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), RequestId("requestId"), GenerateID())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func WithTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(RequestId("requestId")).(string)
}

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"*",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	})
}

type stackError struct {
	stack []byte
	err   error
}

func NewStackError(err error) *stackError {
	var buf [16 * 1024]byte
	n := runtime.Stack(buf[:], false)
	return &stackError{
		stack: buf[:n],
		err:   err,
	}
}

func (e stackError) Error() string {
	return e.err.Error()
}

func (e *stackError) Unwrap() error {
	return e.err
}

type RequestInfo struct {
	status                                            int
	contents_length                                   int64
	method, path, sourceIP, query, user_agent, errors string
	elapsed                                           time.Duration
}

func RequestLogger(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		slog.Log(c, logger.SeverityInfo, "処理開始", "request Id", GetRequestID(c.Request.Context()))
		c.Next()

		r := &RequestInfo{
			status:          c.Writer.Status(),
			contents_length: c.Request.ContentLength,
			path:            c.Request.URL.Path,
			sourceIP:        c.ClientIP(),
			query:           c.Request.URL.RawQuery,
			user_agent:      c.Request.UserAgent(),
			errors:          c.Errors.ByType(gin.ErrorTypePrivate).String(),
			elapsed:         time.Since(start),
		}
		slog.Log(c, logger.SeverityInfo, "処理終了", "Request", r.LogValue(), "requestId", GetRequestID(c.Request.Context()))
	}
}

func (r *RequestInfo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status", r.status),
		slog.Int64("Content-length", r.contents_length),
		slog.String("method", r.method),
		slog.String("path", r.path),
		slog.String("sourceIP", r.sourceIP),
		slog.String("query", r.query),
		slog.String("user_agent", r.user_agent),
		slog.String("errors", r.errors),
		slog.String("elapsed", r.elapsed.String()),
	)
}

func GenerateID() string {
	return uuid.NewString()
}
