package store

import (
	"context"
	"fmt"

	"github.com/767829413/normal-frame/internal/pkg/options"

	"github.com/go-redis/redis/v8"
)

type myRedis struct {
	client *redis.Client
	prefix string
}

func (r *myRedis) Getclient() *redis.Client {
	return r.client
}

func (r *myRedis) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

var (
	client *myRedis
)

func GetRedisIncOr(opts *options.RedisOptions) *myRedis {
	if client != nil {
		return client
	}
	if opts != nil && !opts.Enabled {
		return nil
	}
	if opts == nil && client == nil {
		return nil
	}
	var err error
	once.Do(func() {
		tmpClient := redis.NewClient(&redis.Options{
			Addr: opts.Address,
		})
		if err = tmpClient.Ping(context.Background()).Err(); err != nil {
			return
		}
		client = &myRedis{client: tmpClient, prefix: opts.Prefix}
	})
	if err != nil {
		panic(fmt.Sprintf("GetRedisIncOr err : %v", err))
	}
	return client
}
