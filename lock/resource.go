package lock

import (
	"concurrent-res/config"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"reflect"
)

type Course struct {
	ID   string
	Name string
}

type LockedResource struct {
	Alias      string
	Field      string
	AccessorID string
	Lock       Lock
	Redis      *redis.Client
}

func (lr *LockedResource) Init(conf *config.RedisConfig, expectedInitial int64) error {
	client, err := config.NewRedisClient(conf)
	rl, err := NewRedis(conf)
	if err != nil {
		return err
	}
	lr.Redis = client
	lr.Lock = rl
	_, err = lr.acquireWritingLock(func() (interface{}, error) {
		exists, err := lr.Redis.Exists(lr.Alias).Result()
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			lr.Redis.HSet(lr.Alias, lr.Field, expectedInitial)
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (lr LockedResource) GetRemaining() (int64, error) {
	cmd := lr.Redis.HGet(lr.Alias, lr.Field)
	return cmd.Int64()
}

func (lr LockedResource) acquireWritingLock(op func() (interface{}, error)) (interface{}, error) {
	var lockKey = fmt.Sprintf("%s_persisting", lr.Alias)
	refCount, err := TryLock(lr.Lock, lockKey, lr.AccessorID, false, 0, 5)
	if err != nil {
		return nil, err
	}
	if refCount == 0 {
		return nil, errors.New("res is busy")
	}
	defer lr.Lock.Unlock(lockKey, lr.AccessorID)
	return op()
}

func (lr LockedResource) SetRemainingBy(delta int64, latch int64) (int64, error) {
	result, err := lr.acquireWritingLock(func() (interface{}, error) {
		rem, err := lr.Redis.HGet(lr.Alias, lr.Field).Int64()
		if err != nil {
			return latch, err
		}
		rem += delta
		if rem < latch {
			return latch, errors.New("no rem")
		}
		cmd := lr.Redis.HSet(lr.Alias, lr.Field, rem)
		err = cmd.Err()
		if err != nil {
			return latch, err
		}
		return rem, nil
	})
	if err != nil {
		return latch, err
	}
	return reflect.ValueOf(result).Int(), nil
}

func (lr LockedResource) Persistence(update func(newValue int64) error) error {
	//先在缓存中更新，最后把缓存的结果写到数据库里
	//这里用Lock获取锁，如果失败就不由这个线程更新
	//数据库写入出错不影响后续对缓存的读写，只保证最终一致性
	_, err := lr.acquireWritingLock(func() (interface{}, error) {
		rem, err := lr.GetRemaining()
		if err != nil {
			return nil, err
		}
		return nil, update(rem)
	})
	return err
}
