package middleware

import (
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if lrw.wroteHeader {
		return
	}
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
	lrw.wroteHeader = true
}

func (lrw *loggingResponseWriter) Status() int {
	return lrw.statusCode
}

func GetRequestLogFunc(l *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// задача Handler - обработка запросов, поэтому Middleware должен вернуть handler, мы используем HandlerFunc для простоты
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			// Wrap the original ResponseWriter
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// главная функция для продолжения работы, без нее ваш Middleware поломает ответы на запросы, поскольку не передаст управление функциям из Router
			next.ServeHTTP(lrw, r)

			// Log details after the response is written
			duration := time.Since(start).String()

			requester := ""
			var err error
			if requester, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
				l.Warnf("unable to split host port from remote addr: %v", err)
			}
			l.Infof("%s %s %s %d %s", requester, r.Method, r.URL.Path, lrw.Status(), duration)
		})
	}
}
