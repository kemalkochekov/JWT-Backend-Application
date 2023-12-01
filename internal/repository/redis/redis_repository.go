package redis

import (
	"Fiber_JWT_Authentication_backend_server/internal/connectionRedis"
	"Fiber_JWT_Authentication_backend_server/internal/repository/databaseModel"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

type clientRedisRepo struct {
	db *redis.Client
}

func NewClientRedisRepository(db *redis.Client) connectionRedis.CacheRepository {
	return &clientRedisRepo{db: db}
}

var SessionPrefix = "session"
var LockRequestsPrefix = "lock-req"

func (r *clientRedisRepo) getDbKey(sessionKey, userAgent string) string {
	return strings.Join([]string{
		SessionPrefix,
		sessionKey,
		userAgent,
	}, ":")
}

// used for checking in middleware
func (r *clientRedisRepo) getDbKeyByAuthHeaders(args databaseModel.AuthHeaders) string {
	return r.getDbKey(args.SessionKey, args.UserAgent)
}

// used for login
func (r *clientRedisRepo) getDbKeyByCacheClientSession(args databaseModel.CacheUserSession) string {
	return r.getDbKey(args.SessionKey, args.UserAgent)
}

func (r *clientRedisRepo) GetSession(ctx context.Context, authHeaders databaseModel.AuthHeaders) (databaseModel.CacheUserSession, error) {
	result := databaseModel.CacheUserSession{}
	key := r.getDbKeyByAuthHeaders(authHeaders)
	valueString, err := r.db.Get(ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return result, errors.New("Not Authorized User")
	} else if err != nil {
		return result, errors.Wrapf(err, "User.RedisRepo.GetSession.Get()")
	}

	err = json.Unmarshal([]byte(valueString), &result)
	if err != nil {
		return result, errors.Wrapf(err, "User.RedisRepo.GetSession.Unmarshal()")
	}

	return result, nil
}
func (r *clientRedisRepo) PutSession(ctx context.Context, params databaseModel.CacheUserSession) error {
	sessionBytes, err := json.Marshal(params)
	if err != nil {
		return errors.Wrapf(err, "client.repository.PutSession.Marshal()")
	}

	key := r.getDbKeyByCacheClientSession(params)

	_, err = r.db.Set(ctx, key, sessionBytes, params.Duration).Result()
	if err != nil {
		return errors.Wrapf(err, "User.repository.PutSession")
	}

	return nil
}

func (r *clientRedisRepo) DelSession(ctx context.Context, sessionKey string) error {
	regex := r.getDbKey(sessionKey, "*")

	return errors.Wrapf(r.deleteSessionsByRegex(ctx, regex), "User.RedisRepo.DelSession")
}

func (r *clientRedisRepo) DelAllSessions(ctx context.Context, userId int) error {
	regex := r.getDbKey(strconv.Itoa(userId)+":*", "*")

	return errors.Wrapf(r.deleteSessionsByRegex(ctx, regex), "User.RedisRepo.DelAllSessions")
}

func (r *clientRedisRepo) deleteSessionsByRegex(ctx context.Context, regex string) error {
	keys := make([]string, 0, 10)

	iter := r.db.Scan(ctx, 0, regex, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if len(keys) == 0 {
		return nil
	}

	_, err := r.db.Del(ctx, keys...).Result()
	if err != nil {
		return errors.Wrapf(err, "User.repository.FlushClientSessions.Del(%+v)", keys)
	}

	return nil
}

func (r *clientRedisRepo) LockRequests(ctx context.Context, phoneNumber string) (bool, error) {
	set, err := r.db.SetNX(
		ctx,
		fmt.Sprintf("%s:%s", LockRequestsPrefix, phoneNumber),
		fmt.Sprintf("%d", time.Now().Unix()),
		1*time.Second,
	).Result()
	if err != nil {
		return set, err
	}

	return set, nil
}
