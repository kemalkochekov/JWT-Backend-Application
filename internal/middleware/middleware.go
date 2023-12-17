package middleware

import (
	"Fiber_JWT_Authentication_backend_server/pkg/connectionRedis"
)

type MDWManager struct {
	officiantRedisRepo connectionRedis.CacheRepository
}

func NewOfficiantMiddleware(officiantRedisRepo connectionRedis.CacheRepository) *MDWManager {
	return &MDWManager{
		officiantRedisRepo: officiantRedisRepo,
	}
}
