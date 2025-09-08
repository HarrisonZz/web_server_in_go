package deps

import (
	"github.com/HarrisonZz/web_server_in_go/internal/cache"
	"github.com/gin-gonic/gin"
)

const ctxKeyCache = "cache"

type Deps struct {
	Cache cache.Cache
}

func InjectDeps(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ctxKeyCache, d.Cache)
		c.Next()
	}
}

func CacheFrom(c *gin.Context) cache.Cache {
	v, ok := c.Get(ctxKeyCache)
	if !ok {
		return nil
	}
	if cc, ok := v.(cache.Cache); ok {
		return cc
	}
	return nil
}