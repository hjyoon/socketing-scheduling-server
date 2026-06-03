package schedule

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	c *redis.Client
}

func NewRedis(c *redis.Client) *Redis { return &Redis{c: c} }

func (r *Redis) Ready(ctx context.Context) error { return r.c.Ping(ctx).Err() }

func (r *Redis) Publish(ctx context.Context, channel string, message any) error {
	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return r.c.Publish(ctx, channel, raw).Err()
}

func (r *Redis) QueueLength(ctx context.Context, name string) (int, error) {
	n, err := r.c.LLen(ctx, name).Result()
	return int(n), err
}

func (r *Redis) RoomCount(ctx context.Context, room string) (int, error) {
	n, err := r.c.Get(ctx, "room:"+room+":count").Int()
	if err == redis.Nil {
		return 0, nil
	}
	return n, err
}
