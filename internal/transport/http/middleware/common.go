package http_middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	auth_service "github.com/Phanile/go-secure-labs/internal/features/auth/service"
	"github.com/Phanile/go-secure-labs/internal/logger"
	http_response "github.com/Phanile/go-secure-labs/internal/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIdHeader     = "X-Request-ID"
	logKeyForContext    = "log"
	cookieSessionKey    = "session_id"
	userIdKeyForContext = "user_id"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHeader)
			if requestId == "" {
				requestId = uuid.NewString()
			}

			r.Header.Set(requestIdHeader, requestId)
			w.Header().Set(requestIdHeader, requestId)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(l *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Logger.With(
				zap.String("request_id", r.Header.Get(requestIdHeader)),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), logKeyForContext, l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			responseHandler := http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "Panic() unexpected panic by http request")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)

			beforeTime := time.Now()

			log.Debug(
				fmt.Sprintf("incoming %s HTTP request %s", r.Method, r.Header.Get(requestIdHeader)),
				zap.Time("time", time.Now().UTC()),
			)

			next.ServeHTTP(w, r)

			log.Debug(
				fmt.Sprintf("HTTP request %s is done: ", r.Header.Get(requestIdHeader)),
				zap.Duration("latency", time.Since(beforeTime)),
			)
		})
	}
}

func Auth(a *auth_service.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			sessionIdFromCookie, err := r.Cookie(cookieSessionKey)
			if err != nil || sessionIdFromCookie.Value == "" {
				next.ServeHTTP(w, r)
				return
			}

			userId, err := a.VerifySession(ctx, sessionIdFromCookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx = context.WithValue(ctx, userIdKeyForContext, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(a *auth_service.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(r.Context())
			responseHandler := http_response.NewHTTPResponseHandler(log, w)

			sessionIdFromCookie, err := r.Cookie(cookieSessionKey)
			if err != nil || sessionIdFromCookie.Value == "" {
				responseHandler.ErrorResponse(err, "missing session cookie")
				return
			}

			userId, err := a.VerifySession(ctx, sessionIdFromCookie.Value)
			if err != nil {
				responseHandler.ErrorResponse(err, "invalid session")
				return
			}

			ctx = context.WithValue(ctx, userIdKeyForContext, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIdFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIdKeyForContext).(int64)
	return id, ok
}
