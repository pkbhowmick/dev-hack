package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func AddTracing(tracerName, spanName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		prop := otel.GetTextMapPropagator()
		r := c.Request
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

		newReq := c.Request.WithContext(ctx)

		c.Request = newReq

		c.Next()

		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(c.Writer.Status()),
		)
	}
}
