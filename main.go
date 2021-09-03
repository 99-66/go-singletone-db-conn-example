package main

import (
	"context"
	"github.com/99-66/go-singletone-db-conn-example/repositories"
	"log"
)

func init() {
	repositories.Init()
}

func main() {
	// MySQL: Close the connection when the main function exits
	mysql := repositories.GetMySQL()
	if mysql != nil {
		db, _ := mysql.DB()
		defer db.Close()
	}

	// PostgreSQL: Close the connection when the main function exits
	postgresql := repositories.GetPostgreSQL()
	if postgresql != nil {
		db, _ := postgresql.DB()
		defer db.Close()
	}

	// MongoDB : Close the connection when the main function exits
	mongodb := repositories.GetMongoDB()
	if mongodb != nil {
		defer mongodb.Disconnect(context.Background())
	}

	// ElasticSearch
	es := repositories.GetElasticSearch()
	esInfo, err := es.Info()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n", esInfo)

	// Redis: Close the connection when the main function exits
	redis := repositories.GetRedis()
	if redis != nil {
		defer redis.Close()
	}
	status := redis.Ping(context.TODO())
	log.Printf("%+v\n", status.String())

	// Kafka(producer):	Close the connection when the main function exits
	producer := repositories.GetKafkaProducer()
	if producer != nil {
		defer producer.Close()
	}

	// Kafka(consumer): Close the connection when the main function exits
	consumer := repositories.GetKafkaConsumer()
	if consumer != nil {
		defer consumer.Close()
	}
}
