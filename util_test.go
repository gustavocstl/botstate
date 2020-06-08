package botstate_test

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v7"
	"github.com/gucastiliao/botstate"
)

func mockRedis() {
	mr, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	r := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	botstate.SetStorageClient(botstate.DefaultStorage(r))
}
