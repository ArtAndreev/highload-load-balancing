package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	AccessHits *prometheus.CounterVec
)

func InitMetrics(ns string) {
	AccessHits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Name:      "hits_by_http_status",
		Help:      "Total hits ordered by http response statuses",
	},
		[]string{"http_status", "path", "method"},
	)
}

func CountHitsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &StatusWrapperForResponseWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(ww, r)

		AccessHits.With(prometheus.Labels{
			"http_status": strconv.Itoa(ww.Status),
			"path":        r.URL.Path,
			"method":      r.Method,
		}).Inc()
	})
}
