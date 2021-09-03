package repositories

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	kafkaConsumerOnce   sync.Once
	kafkaConsumer       *cluster.Consumer
	kafkaConsumerConfig consumerConfig
)

type consumerConfig struct {
	Topic         []string `mapstructure:"KAFKA_CONSUMER_TOPIC"`
	Brokers       []string `mapstructure:"KAFKA_CONSUMER_BROKER"`
	ConsumerGroup string   `mapstructure:"KAFKA_CONSUMER_GROUP"`
}

func kafkaConnConfig() *cluster.Config {
	conf := cluster.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Group.Return.Notifications = true
	conf.Consumer.Offsets.CommitInterval = time.Second

	return conf
}

// newKafkaConsumer Kafka Consumer 생성자
// 이 함수는 singleton pattern 을 위해 once.Do()에서 실행할 수 있게 함수로 반환한다
func newKafkaConsumer() func() {
	var err error
	return func() {
		kafkaConsumer, err = cluster.NewConsumer(
			kafkaConsumerConfig.Brokers,
			kafkaConsumerConfig.ConsumerGroup,
			kafkaConsumerConfig.Topic,
			kafkaConnConfig(),
		)
		if err != nil {
			panic(err)
		}
	}
}

// init 환경 변수를 불러온다
func init() {
	LoadConfig(".")

	err := viper.Unmarshal(&kafkaConsumerConfig)
	if err != nil {
		panic(err)
	}
}

// initKafkaConsumer init.go 에서 kafka consumer connection 초기화 시에 사용하는 함수이다
func initKafkaConsumer() {
	if kafkaConsumer == nil {
		kafkaConsumerOnce.Do(newKafkaConsumer())
	}
}

// GetKafkaConsumer Kafka Consumer 객체를 반환하는 함수이다
func GetKafkaConsumer() *cluster.Consumer {
	if kafkaConsumer == nil {
		kafkaConsumerOnce.Do(newKafkaConsumer())
	}

	return kafkaConsumer
}
