package databaseModel

import "time"

type AuthHeaders struct {
	UserAgent  string
	SessionKey string
}

type CacheUserSession struct {
	// fields for keys in redis
	SessionKey string // userId + random uuid
	UserAgent  string `json:"userAgent"`
	Duration   time.Duration
	CreatedAt  int64 `json:"createdAt"` // unix milli
}
