package repositories

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"sync"
)

var (
	kafkaProducerOnce   sync.Once
	kafkaProducer       sarama.AsyncProducer
	kafkaProducerConfig producerConfig
)

type producerConfig struct {
	Topic   string   `mapstructure:"KAFKA_PRODUCER_TOPIC"`
	Brokers []string `mapstructure:"KAFKA_PRODUCER_BROKER"`
}

func producerConnConfig() *sarama.Config {
	return sarama.NewConfig()
}

// newKafkaProducer Kafka Producer 생성자
// 이 함수는 singleton pattern 을 위해 once.Do()에서 실행할 수 있게 함수로 반환한다
func newKafkaProducer() func() {
	var err error
	return func() {
		config := producerConnConfig()
		kafkaProducer, err = sarama.NewAsyncProducer(kafkaProducerConfig.Brokers, config)
		if err != nil {
			panic(err)
		}
	}
}

// init 환경변수를 불러온다
func init() {
	LoadConfig(".")

	err := viper.Unmarshal(&kafkaProducerConfig)
	if err != nil {
		panic(err)
	}
}

// initKafkaProducer init.go 에서 kafka producer connection 초기화 시에 사용하는 함수이다
func initKafkaProducer() {
	if kafkaProducer == nil {
		kafkaProducerOnce.Do(newKafkaProducer())
	}
}

// GetKafkaProducer Kafka Producer 객체를 반환하는 함수이다
func GetKafkaProducer() sarama.AsyncProducer {
	if kafkaProducer == nil {
		kafkaProducerOnce.Do(newKafkaProducer())
	}

	return kafkaProducer
}
