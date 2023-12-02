package connectionRedis

import (
	"Fiber_JWT_Authentication_backend_server/internal/repository/databaseModel"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	GetSession(ctx context.Context, sessionKey databaseModel.AuthHeaders) (databaseModel.CacheUserSession, error)
	PutSession(ctx context.Context, params databaseModel.CacheUserSession) error
	DelSession(ctx context.Context, sessionKey string) error
	DelAllSessions(ctx context.Context, clientUserId int) error
	LockRequests(ctx context.Context, username string) (bool, error)
}

type Database struct {
	Client *redis.Client
}

var (
	ErrNil = errors.New("no matching record found in redis database")
	Ctx    = context.TODO()
)

func NewDatabase(ctx context.Context) (*Database, error) {
	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, err
	}
	return &Database{
		Client: client,
	}, nil
}
