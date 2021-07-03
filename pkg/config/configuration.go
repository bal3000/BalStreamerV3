package config

type Configuration struct {
	ConnectionString string `json:"connectionString"`
	RabbitURL        string `json:"rabbitUrl"`
	ExchangeName     string `json:"exchangeName"`
	QueueName        string `json:"queueName"`
	Durable          bool   `json:"durable"`
	LiveStreamURL    string `json:"liveStreamUrl"`
	APIKey           string `json:"apiKey"`
	CasterURL        string `json:"casterUrl"`
}
