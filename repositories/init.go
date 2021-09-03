package repositories

func Init() {
	// Initializing Mysql Connection
	initMySQL()
	// Initializing PostgreSQL Connection
	initPostgreSQL()
	// Initializing MongoDB Connection
	initMongoDB()
	// Initializing ElasticSearch Connection
	initElasticSearch()
	// Initializing Kafka Producer Connection
	initKafkaProducer()
	// Initializing Kafka Consumer Connection
	initKafkaConsumer()
}
