package lock

import (
	"concurrent-res/config"
	"github.com/go-redis/redis"
	"time"
)

type Redis struct {
	Client        *redis.Client
	LockScript    string
	ReleaseScript string
	Prefix        string
}

func NewRedis(conf *config.RedisConfig) (*Redis, error) {
	client, err := config.NewRedisClient(conf)
	if err != nil {
		return nil, err
	}
	return &Redis{
		Client:        client,
		LockScript:    LOCK_LUA,
		ReleaseScript: RELEASE_LUA,
		Prefix:        "lock",
	}, nil
}

func (lock Redis) Lock(alias string, owner string, repeat bool, expireSec int) (int, error) {
	var key = lock.Prefix + "_" + alias
	var repeatVal = 0
	if repeat {
		repeatVal = 1
	}
	var start = time.Now()
	var cmd = lock.Client.Eval(lock.LockScript, []string{key}, owner, repeatVal, expireSec)
	var netSpent int = int(time.Now().Sub(start) / time.Second)
	result, err := cmd.Int()
	if err != nil {
		return result, err
	}
	if expireSec < netSpent { //远程锁实际已经超时过期
		return 0, nil
	}
	return result, err
}

func (lock Redis) Unlock(alias string, owner string) (int, error) {
	var key = lock.Prefix + "_" + alias
	var cmd = lock.Client.Eval(lock.ReleaseScript, []string{key}, owner)
	return cmd.Int()
}
