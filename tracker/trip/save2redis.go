package trip

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var (
	RedisServer         string
	RedisServerPassword string
	RedisClient         *redis.Client
)

func init() {
	RedisServer = os.Getenv("REDIS_SERVER_ADDRESS")
	RedisServerPassword = os.Getenv("REDIS_SERVER_PASSWORD")
	if RedisServer == "" {
		log.Warn().Str("REDIS_SERVER_ADDRESS", RedisServer).
			Msg("No value for REDIS_SERVER_ADDRESS and failed to connect Redis server.")

	} else {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     RedisServer,
			Password: RedisServerPassword,
		})
		// re-try connection util 120s
		timeout := time.Now().Add(120 * time.Second)
		for {
			pong, err := RedisClient.Ping(context.TODO()).Result()
			if err != nil {
				log.Error().Str("to", RedisServer).Str("result", pong).Msg("Failed to connect Redis.")
			} else {
				break
			}
			if timeout.After(time.Now()) {
				break
			}

		}
	}
}

func SaveTd2Redis(id string, buf []byte) error {

	pipeline := RedisClient.Pipeline()
	ctx := context.TODO()
	pipeline.HSet(ctx, "trip_detail_cache", id, buf)
	cmds, err := pipeline.Exec(ctx)

	for _, cmd := range cmds {
		log.Info().Interface("args", cmd.Args()).Str("cmd", cmd.FullName()).Msg("Execute commands")
	}

	return err
}

func QueryTd4Redis(id string) ([]byte, error) {
	ctx := context.TODO()
	return RedisClient.HGet(ctx, "trip_detail_cache", id).Bytes()

}

func QueryAllTds4Redis() (map[string]string, error) {
	ctx := context.TODO()
	return RedisClient.HGetAll(ctx, "trip_detail_cache").Result()

}

func ClearHistory() {
	ctx := context.TODO()
	keys, _ := RedisClient.HKeys(ctx, "trip_detail_cache").Result()
	for _, k := range keys {
		cmd := RedisClient.HDel(ctx, "trip_detail_cache", k)
		log.Debug().Str("cmd", cmd.FullName()).Interface("args", cmd.Args()).Msg("ClearHistory()")
	}
}
