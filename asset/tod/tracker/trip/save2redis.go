package trip

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type S2Redis struct {
	Server   string
	Password string
	*redis.Client
}

type S2RedisInterface interface {
	Connect() error
	SaveTripDetail(id string, buf []byte) error
	TripDetail(id string) ([]byte, error)
	AllTripDetail() (map[string]string, error)
	ClearTripDetail() error
}

func (s2r *S2Redis) Connect() error {
	if s2r.Client == nil {
		if s2r.Server == "" {
			s2r.Server = "127.0.0.1:6379"
			log.Warn().Msg("No value for Redis Server, connect to " + s2r.Server)
		}
		s2r.Client = redis.NewClient(&redis.Options{
			Addr:     s2r.Server,
			Password: s2r.Password,
		})
		// re-try connection util 120s

		go func() {
			timeout := time.Now().Add(120 * time.Second)
			for {
				pong, err := s2r.Ping(context.TODO()).Result()
				if err != nil {
					log.Error().Str("to", s2r.Server).Str("result", pong).Msg("Failed to connect Redis.")
				} else {
					break
				}
				if timeout.After(time.Now()) {
					break
				}
			}
		}()
	}
	return nil
}

func (s2r *S2Redis) SaveTripDetail(id string, buf []byte) error {
	pipeline := s2r.Pipeline()
	ctx := context.TODO()
	pipeline.HSet(ctx, "trip_detail_cache", id, buf)
	cmds, err := pipeline.Exec(ctx)

	for _, cmd := range cmds {
		log.Info().Interface("args", cmd.Args()).Str("cmd", cmd.FullName()).Msg("Execute commands")
	}

	return err
}

func (s2r *S2Redis) TripDetail(id string) ([]byte, error) {
	ctx := context.TODO()
	return s2r.HGet(ctx, "trip_detail_cache", id).Bytes()
}

func (s2r *S2Redis) AllTripDetail() (map[string]string, error) {
	ctx := context.TODO()
	return s2r.HGetAll(ctx, "trip_detail_cache").Result()
}

func (s2r *S2Redis) ClearTripDetail() error {
	ctx := context.TODO()
	keys, _ := s2r.HKeys(ctx, "trip_detail_cache").Result()
	for _, k := range keys {
		cmd := s2r.HDel(ctx, "trip_detail_cache", k)
		log.Debug().Str("cmd", cmd.FullName()).Interface("args", cmd.Args()).Msg("ClearHistory()")
	}
	return nil
}

// var (
// 	RedisServer         string
// 	RedisServerPassword string
// 	RedisClient         *redis.Client
// )

// func init() {
// 	RedisServer = os.Getenv("REDIS_SERVER_ADDRESS")
// 	RedisServerPassword = os.Getenv("REDIS_SERVER_PASSWORD")
// 	if RedisServer == "" {
// 		RedisServer = "127.0.0.1:6379"
// 		log.Warn().Str("REDIS_SERVER_ADDRESS", RedisServer).Msg("No value for REDIS_SERVER_ADDRESS, using default: 127.0.0.1:6379")
// 	}
// 	RedisClient = redis.NewClient(&redis.Options{
// 		Addr:     RedisServer,
// 		Password: RedisServerPassword,
// 	})
// 	// re-try connection util 120s

// 	go func() {
// 		timeout := time.Now().Add(120 * time.Second)
// 		for {
// 			pong, err := RedisClient.Ping(context.TODO()).Result()
// 			if err != nil {
// 				log.Error().Str("to", RedisServer).Str("result", pong).Msg("Failed to connect Redis.")
// 			} else {
// 				break
// 			}
// 			if timeout.After(time.Now()) {
// 				break
// 			}
// 		}
// 	}()

// }

// func SaveTd2Redis(id string, buf []byte) error {

// 	pipeline := RedisClient.Pipeline()
// 	ctx := context.TODO()
// 	pipeline.HSet(ctx, "trip_detail_cache", id, buf)
// 	cmds, err := pipeline.Exec(ctx)

// 	for _, cmd := range cmds {
// 		log.Info().Interface("args", cmd.Args()).Str("cmd", cmd.FullName()).Msg("Execute commands")
// 	}

// 	return err
// }

// func QueryTd4Redis(id string) ([]byte, error) {
// 	ctx := context.TODO()
// 	return RedisClient.HGet(ctx, "trip_detail_cache", id).Bytes()

// }

// func QueryAllTds4Redis() (map[string]string, error) {
// 	ctx := context.TODO()
// 	return RedisClient.HGetAll(ctx, "trip_detail_cache").Result()

// }

// func ClearHistory() {
// 	ctx := context.TODO()
// 	keys, _ := RedisClient.HKeys(ctx, "trip_detail_cache").Result()
// 	for _, k := range keys {
// 		cmd := RedisClient.HDel(ctx, "trip_detail_cache", k)
// 		log.Debug().Str("cmd", cmd.FullName()).Interface("args", cmd.Args()).Msg("ClearHistory()")
// 	}
// }
