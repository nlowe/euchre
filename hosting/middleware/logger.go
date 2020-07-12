package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func LogrusLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}

				fields := logrus.Fields{
					"prefix":   "http",
					"code":     ww.Status(),
					"len":      ww.BytesWritten(),
					"duration": time.Since(start).String(),
				}

				logrus.WithFields(fields).Infof("%s %s://%s%s", r.Method, scheme, r.Host, r.RequestURI)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
