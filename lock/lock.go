package lock

import "time"

type Lock interface {
	//TryLock(alias string,owner string,repeat bool,expireSec int,timeOutSec int) (int,error)
	Lock(alias string, owner string, repeat bool, expireSec int) (int, error)
	Unlock(alias string, owner string) (int, error)
}

func TryLock(lock Lock, alias string, owner string, repeat bool, expireSec int, timeOutSec int) (int, error) {
	var start = time.Now()
	for {
		result, err := lock.Lock(alias, owner, repeat, expireSec)
		if err != nil || result != 0 {
			return result, err
		}
		var timeOutTime = start.Add(time.Duration(timeOutSec) * time.Second)
		if time.Now().After(timeOutTime) {
			return 0, err
		}
	}
}
