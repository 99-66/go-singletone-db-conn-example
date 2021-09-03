package repositories

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"sync"
)

var (
	redisOnce sync.Once
	r         *redis.Client
	rConfig   redisConfig
)

type redisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	DB       int    `mapstructure:"REDIS_DB"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

func (r *redisConfig) dsn() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r *redisConfig) options() *redis.Options {
	return &redis.Options{
		Addr:     r.dsn(),
		Password: r.Password,
		DB:       r.DB,
	}
}

// newRedis Redis connection 생성자
// 이 함수는 singleton pattern 을 위해 once.Do()에서 실행할 수 있게 함수로 반환한다
func newRedis() func() {
	return func() {
		r = redis.NewClient(rConfig.options())
	}
}

// init 환경변수를 불러온다
func init() {
	LoadConfig(".")

	err := viper.Unmarshal(&rConfig)
	if err != nil {
		panic(err)
	}
}

// initRedis init.go 에서 database connection 초기화 시에 사용하는 함수이다.
func initRedis() {
	if r == nil {
		redisOnce.Do(newRedis())
	}
}

// GetRedis Redis 객체를 반환하는 함수이다
func GetRedis() *redis.Client {
	if r == nil {
		redisOnce.Do(newRedis())
	}

	return r
}
