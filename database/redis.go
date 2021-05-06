package database

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/fabjan/robotoscope/core"
	"github.com/go-redis/redis/v8"
)

// RedisStore tracks robots in Redis.
type RedisStore struct {
	client *redis.Client
	prefix string
}

func (rs *RedisStore) key(name string) string {
	return rs.prefix + name
}

// Count increases the seen count for the given bot.
func (rs *RedisStore) Count(name string) error {
	ctx := context.Background()
	return rs.client.Incr(ctx, rs.key(name)).Err()
}

// List returns a list showing how many times each robot has been seen.
func (rs *RedisStore) List() ([]core.RobotInfo, error) {
	info := []core.RobotInfo{}

	ctx := context.Background()
	bots, err := rs.client.Keys(ctx, rs.key("*")).Result()
	if err != nil {
		return nil, err
	}

	if len(bots) == 0 {
		return info, nil
	}

	counts, err := rs.client.MGet(ctx, bots...).Result()
	if err != nil {
		return nil, err
	}

	for i, agent := range bots {
		ri := core.RobotInfo{
			UserAgent: strings.TrimPrefix(agent, rs.prefix),
		}
		switch c := counts[i].(type) {
		case int:
			ri.Seen = c
		case string:
			s, err := strconv.Atoi(c)
			if err != nil {
				return nil, errors.New("non-integer 'seen' count")
			}
			ri.Seen = s
		default:
			return nil, errors.New("non-integer 'seen' count")
		}
		info = append(info, ri)
	}

	return info, nil
}

// OpenRedis creates a client connected to the Redis with the given URL.
func OpenRedis(rawURL string) *redis.Client {
	redisURL, _ := url.Parse(rawURL)
	redisPassword, _ := redisURL.User.Password()
	redisDB := 0
	redisOptions := redis.Options{
		Addr:     redisURL.Host,
		Password: redisPassword,
		DB:       redisDB,
	}
	return redis.NewClient(&redisOptions)
}

// NewRedisStore creates a RedisStore backed by the given Redis client,
// using the given prefix for its counters.
func NewRedisStore(c *redis.Client, name string) *RedisStore {
	return &RedisStore{
		client: c,
		prefix: name + "::",
	}
}
