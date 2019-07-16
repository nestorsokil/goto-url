package util

import (
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func HttpSpan(spanName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan(spanName)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(r.Context(), span)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
