package repositories

import (
	"github.com/cenkalti/backoff/v4"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	esOnce   sync.Once
	es       *elastic.Client
	esConfig elasticSearchConfig
)

type elasticSearchConfig struct {
	Host     []string `mapstructure:"ELS_HOST"`
	User     string   `mapstructure:"ELS_USER"`
	Password string   `mapstructure:"ELS_PASSWORD"`
}

func (e *elasticSearchConfig) config() *elastic.Config {
	retryBackoff := backoff.NewExponentialBackOff()

	return &elastic.Config{
		Addresses:     e.Host,
		Username:      e.User,
		Password:      e.Password,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	}
}

func newElasticSearch() func() {
	var err error
	return func() {
		es, err = elastic.NewClient(*esConfig.config())
		if err != nil {
			panic(err)
		}
	}
}

// init 환경변수를 불러온다
func init() {
	LoadConfig(".")

	err := viper.Unmarshal(&esConfig)
	if err != nil {
		panic(err)
	}
}

func initElasticSearch() {
	if es == nil {
		esOnce.Do(newElasticSearch())
	}
}

// GetElasticSearch ElasticSearch 객체를 반환하는 함수이다
func GetElasticSearch() *elastic.Client {
	if es == nil {
		esOnce.Do(newElasticSearch())
	}

	return es
}
