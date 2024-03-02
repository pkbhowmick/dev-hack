package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func AddTracing(tracerName, spanName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			prop := otel.GetTextMapPropagator()
			propCtx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := otel.Tracer(tracerName).Start(propCtx, spanName)
			defer span.End()

			span.SetAttributes(
				semconv.HTTPMethodKey.String(r.Method),
				semconv.HTTPRouteKey.String(r.URL.Path),
				semconv.HTTPURLKey.String(r.URL.String()),
				semconv.HTTPHostKey.String(r.Host),
				semconv.HTTPSchemeKey.String(r.URL.Scheme),
			)

			respWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(respWriter, r.WithContext(ctx))

			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int(respWriter.Status()),
			)
		}
		return http.HandlerFunc(fn)
	}
}
