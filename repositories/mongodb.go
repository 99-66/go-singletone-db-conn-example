package repositories

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	mongodbOnce   sync.Once
	mongoDB       *mongo.Client
	mongodbConfig MongoDBConfig
)

type MongoDBConfig struct {
	Host     []string `mapstructure:"MONGODB_HOST"`
	Port     int      `mapstructure:"MONGODB_PORT"`
	DB       string   `mapstructure:"MONGODB_DB"`
	User     string   `mapstructure:"MONGODB_USER"`
	Password string   `mapstructure:"MONGODB_PASSWORD"`
	SSL      bool     `mapstructure:"MONGODB_SSL"`
}

func (m *MongoDBConfig) hosts() (h []string) {
	for _, host := range m.Host {
		h = append(h, fmt.Sprintf("%s:%d", host, m.Port))
	}

	return h
}

func (m *MongoDBConfig) dsn() string {
	return fmt.Sprintf("mongodb://%s:%s@%s/?ssl=%t",
		m.User,
		m.Password,
		m.hosts(),
		m.SSL)
}

func (m *MongoDBConfig) mongoDbOption() *options.ClientOptions {
	return options.Client().ApplyURI(m.dsn())
}

// init 환경변수를 불러온다
func init() {
	LoadConfig(".")

	err := viper.Unmarshal(&mongodbConfig)
	if err != nil {
		panic(err)
	}
}

// newMongoDB MongoDB connection 생성자
// 이 함수는 singleton pattern 을 위해 once.Do()에서 실행할 수 있게 함수로 반환한다
func newMongoDB() func() {
	var err error
	return func() {
		dbOptions := mongodbConfig.mongoDbOption()
		mongoDB, err = mongo.Connect(context.Background(), dbOptions)
		if err != nil {
			panic(err)
		}
	}
}

// initMongoDB init.go 에서 database connection 초기화 시에 사용하는 함수이다.
func initMongoDB() {
	if mongoDB == nil {
		mongodbOnce.Do(newMongoDB())
	}
}

// GetMongoDB MongoDB 객체를 반환하는 함수이다
func GetMongoDB() *mongo.Client {
	if mongoDB == nil {
		mongodbOnce.Do(newMongoDB())
	}

	return mongoDB
}
