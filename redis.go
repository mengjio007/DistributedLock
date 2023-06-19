package DistributedLock

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

var (
	lockScript = NewScript(`if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`)
	delScript = NewScript(`if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`)
)

// NewScript returns a new Script instance.
func NewScript(script string) *redis.Script {
	return redis.NewScript(script)
}

type RedisLock struct {
	store   *redis.Client
	key     string
	id      string
	timeout string
}

func NewRedisLock(client *redis.Client, key, id, timeout string) Lock {
	if len(id) <= 0 {
		id = time.Now().String()
	}
	return &RedisLock{
		store:   client,
		key:     key,
		id:      id,
		timeout: timeout,
	}
}

func (r *RedisLock) Lock() error {
	result, err := lockScript.Run(r.store, []string{r.key}, r.id, r.timeout).Result()
	if err != nil {
		return err
	}
	if reply, ok := result.(string); ok && reply == "OK" {
		return nil
	}
	return errors.New("have some diff error")
}

func (r *RedisLock) UnLock() error {
	resp, err := delScript.Run(r.store, []string{r.key}, r.id).Result()
	if err != nil {
		return err
	}
	if reply, ok := resp.(int64); ok && reply == 1 {
		return nil
	}
	return errors.New("have some diff error")
}

//type RedisRedLock struct {
//	RedisPool []*redis.Client
//}
//
//func NewRedisRedLock() {
//
//}
