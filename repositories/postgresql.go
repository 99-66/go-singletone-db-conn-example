package repositories

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var (
	postgreOnce      sync.Once
	postgreSQL       *gorm.DB
	postgreSqlConfig PostgreSQLConfig
)

type PostgreSQLConfig struct {
	Host     string `mapstructure:"POSTGRESQL_HOST"`
	Port     int    `mapstructure:"POSTGRESQL_PORT"`
	DB       string `mapstructure:"POSTGRESQL_DB"`
	User     string `mapstructure:"POSTGRESQL_USER"`
	Password string `mapstructure:"POSTGRESQL_PASSWORD"`
}

func (p *PostgreSQLConfig) dsn() string {
	timezone := "Asia/Seoul"
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		postgreSqlConfig.Host,
		postgreSqlConfig.User,
		postgreSqlConfig.Password,
		postgreSqlConfig.DB,
		postgreSqlConfig.Port,
		timezone)
}

func postgreSqlDbConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
}

// newPostgreSQL PostgreSQL Connection 생성자
// 이 함수는 singleton pattern 을 위해 once.Do()에서 실행할 수 있게 함수로 반환한다
func newPostgreSQL() func() {
	var err error
	return func() {
		dbConfig := postgreSqlDbConfig()
		postgreSQL, err = gorm.Open(postgres.Open(postgreSqlConfig.dsn()), dbConfig)
		if err != nil {
			panic(err)
		}
	}
}

// init 환경변수를 불러온다
func init() {
	LoadConfig(".")

	err := viper.Unmarshal(&postgreSqlConfig)
	if err != nil {
		panic(err)
	}
}

// initPostgreSQL init.go 에서 database connection 초기화 시에 사용하는 함수이다
func initPostgreSQL() {
	if postgreSQL == nil {
		postgreOnce.Do(newPostgreSQL())
	}
}

// GetPostgreSQL PostgreSQL 객체를 반환하는 함수이다
func GetPostgreSQL() *gorm.DB {
	if postgreSQL == nil {
		postgreOnce.Do(newPostgreSQL())
	}

	return postgreSQL
}
