package main

import (
	"concurrent-res/config"
	"concurrent-res/lock"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	rl, err := lock.NewRedis(&config.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	})
	if err != nil {
		panic(err)
	}
	var val = 0
	for i := 0; i < 30; i++ {
		var threadId = fmt.Sprintf("thread_%d", i)
		go func() {
			for iSub := 0; iSub < 10; iSub++ {
				result, err := lock.TryLock(rl, "c", threadId, true, 30, 30)
				if err != nil {
					panic(err)
				}
				val += 1
				fmt.Println(threadId, "lock:", result, "val:", val)
			}
			for iSub := 0; iSub < 10; iSub++ {
				unlockResult, err := rl.Unlock("c", threadId)
				if err != nil {
					panic(err)
				}
				fmt.Println(threadId, "unlock:", unlockResult)
			}
		}()
	}
	var sigWait = make(chan os.Signal)
	signal.Notify(sigWait, os.Interrupt)
	<-sigWait
}
