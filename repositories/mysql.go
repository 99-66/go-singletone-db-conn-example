package repositories

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var (
	mysqlOnce   sync.Once
	mySql       *gorm.DB
	mysqlConfig MySQLConfig
)

type MySQLConfig struct {
	Host     string `mapstructure:"MYSQL_HOST"`
	Port     int    `mapstructure:"MYSQL_PORT"`
	User     string `mapstructure:"MYSQL_USER"`
	Password string `mapstructure:"MYSQL_PASSWORD"`
}

func (m *MySQLConfig) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port)
}

func mySqlDbConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
}

// newMySQL MySQL connection 생성자
// 이 함수는 singleton pattern 을 위해 once.Do()에서 실행할 수 있게 함수로 반환한다
func newMySQL() func() {
	var err error
	return func() {
		dbConfig := mySqlDbConfig()
		mySql, err = gorm.Open(mysql.Open(mysqlConfig.dsn()), dbConfig)
		if err != nil {
			panic(err)
		}
	}
}

// init 환경변수를 불러온다
func init() {
	LoadConfig(".")

	err := viper.Unmarshal(&mysqlConfig)
	if err != nil {
		panic(err)
	}
}

// initMySQL init.go 에서 database connection 초기화 시에 사용하는 함수이다
func initMySQL() {
	if mySql == nil {
		mysqlOnce.Do(newMySQL())
	}
}

// GetMySQL MySQL 객체를 반환하는 함수이다
func GetMySQL() *gorm.DB {
	if mySql == nil {
		mysqlOnce.Do(newMySQL())
	}

	return mySql
}
