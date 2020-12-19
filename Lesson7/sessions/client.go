package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

//RedisClient struct
type RedisClient struct {
	*redis.Client
	TTL time.Duration
}

//NewRedisClient func
func NewRedisClient(host, port string, ttl time.Duration) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}

	c := &RedisClient{
		TTL:    ttl,
		Client: client,
	}

	return c, nil
}

//Close func
func (c *RedisClient) Close() error {
	return c.Client.Close()
}

//GetRecord func
func (c *RedisClient) GetRecord(mkey string) ([]byte, error) {
	data, err := c.Get(context.Background(), mkey).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return data, nil
}
