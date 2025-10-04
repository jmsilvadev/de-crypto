package config

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	CliUrl                 = getProviderURL()
	DefaultHttpClient      = &http.Client{}
	DefaultRpcUrl          = "https://ethereum-rpc.publicnode.com"
	DefaultCheckpointStore = "./data/checkpoint"
	DefaultAddressFile     = "./data/address.json"

	DefaultHeadsChannelSize  = 64
	DefaultBlocksChannelSize = 64
	DefaultEventsChannelSize = 1024

	DefaultWorkersQty        = 8
	DefaultBatchSize         = 50
	DeaultJitterDeviation    = 0.2
	DefaultMaxEnqueuePerTick = 64

	DefaultPollingInterval = 1 * time.Second
	DefaultRequestTimeout  = 10 * time.Second
	DefaultFlushInterval   = 200 * time.Millisecond

	DefauftKafkaTopic   = "de-crypto-events"
	DefauftKafkaBrokers = []string{"localhost:9092"}
)

type HeadMonitorConfig struct {
	PollInterval time.Duration
	StartFrom    uint64
	Jitter       float64
	// cool thing, I had problems with the buffer blocking the io so I figure it out,
	// we will ignore when the queue is full
	MaxEnqueuePerTick int
}

type BlockFetcherConfig struct {
	Workers        int
	ReqTimeout     time.Duration
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
	Jitter         float64
}

type FilterConfig struct {
	Workers int
}

type SinkConfig struct {
	FlushInterval time.Duration
	BatchSize     int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func GetAddressFile() string {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env not found using fallbacks")
	}

	addressFile := os.Getenv("ADDRESS_FILE")
	if addressFile != "" {
		return addressFile
	}

	return DefaultAddressFile
}

func getProviderURL() string {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env not found using fallbacks")
	}

	providerURL := os.Getenv("RPC_URL")
	if providerURL != "" {
		return providerURL
	}

	return DefaultRpcUrl
}

func GetKafakTopic() string {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env not found using fallbacks")
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic != "" {
		return topic
	}

	return DefauftKafkaTopic
}

func GetKafakBrokers() []string {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env not found using fallbacks")
	}

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers != "" {
		return strings.Split(brokers, ",")
	}

	return DefauftKafkaBrokers
}
