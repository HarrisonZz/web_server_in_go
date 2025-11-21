package telemetry

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func GinChildSpan() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 取得 parent span（由 otelgin middleware 建立）
		tr := otel.Tracer("web-server-in-go/handler")

		// span 名稱使用 path（或你也可以自訂）
		spanName := c.FullPath()
		if spanName == "" {
			spanName = c.Request.URL.Path
		}

		ctx, span := tr.Start(c.Request.Context(), spanName)
		defer span.End()

		// 讓 handler 拿到 span context
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()

		// 交給下一層（真正的 route handler）
		c.Next()

		// 2) 結束後加一些共通 attributes
		span.SetAttributes(
			attribute.String("http.route", c.FullPath()),
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.String("client.ip", c.ClientIP()),
			attribute.Float64("handler.latency_ms", float64(time.Since(start).Milliseconds())),
		)
	}
}
