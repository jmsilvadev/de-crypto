package config

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetProviderURL_WithEnvVariable(t *testing.T) {
	t.Run("get_provider_url_with_rpc_url_env", func(t *testing.T) {
		os.Setenv("RPC_URL", "https://test-rpc-url.com")
		defer os.Unsetenv("RPC_URL")
		url := getProviderURL()
		assert.Equal(t, "https://test-rpc-url.com", url)
	})
	t.Run("get_provider_url_without_env_variable", func(t *testing.T) {
		os.Unsetenv("RPC_URL")
		url := getProviderURL()
		assert.Equal(t, DefaultRpcUrl, url)
	})
	t.Run("get_provider_url_with_empty_env_variable", func(t *testing.T) {
		os.Setenv("RPC_URL", "")
		defer os.Unsetenv("RPC_URL")
		url := getProviderURL()
		assert.Equal(t, DefaultRpcUrl, url)
	})
}

func TestConfigConstants(t *testing.T) {
	t.Run("verify_default_http_client", func(t *testing.T) {
		assert.NotNil(t, DefaultHttpClient)
		assert.IsType(t, &http.Client{}, DefaultHttpClient)
	})
	t.Run("verify_default_rpc_url", func(t *testing.T) {
		assert.Equal(t, "https://ethereum-rpc.publicnode.com", DefaultRpcUrl)
	})
	t.Run("verify_default_checkpoint_store", func(t *testing.T) {
		assert.Equal(t, "./data/checkpoint", DefaultCheckpointStore)
	})
	t.Run("verify_default_address_file", func(t *testing.T) {
		assert.Equal(t, "./data/address.json", DefaultAddressFile)
	})
	t.Run("verify_default_heads_channel_size", func(t *testing.T) {
		assert.Equal(t, 64, DefaultHeadsChannelSize)
	})
	t.Run("verify_default_blocks_channel_size", func(t *testing.T) {
		assert.Equal(t, 64, DefaultBlocksChannelSize)
	})
	t.Run("verify_default_events_channel_size", func(t *testing.T) {
		assert.Equal(t, 1024, DefaultEventsChannelSize)
	})
	t.Run("verify_default_workers_qty", func(t *testing.T) {
		assert.Equal(t, 8, DefaultWorkersQty)
	})
	t.Run("verify_default_batch_size", func(t *testing.T) {
		assert.Equal(t, 50, DefaultBatchSize)
	})
	t.Run("verify_default_jitter_deviation", func(t *testing.T) {
		assert.Equal(t, 0.2, DeaultJitterDeviation)
	})
	t.Run("verify_default_max_enqueue_per_tick", func(t *testing.T) {
		assert.Equal(t, 64, DefaultMaxEnqueuePerTick)
	})
	t.Run("verify_default_polling_interval", func(t *testing.T) {
		assert.Equal(t, 1*time.Second, DefaultPollingInterval)
	})
	t.Run("verify_default_request_timeout", func(t *testing.T) {
		assert.Equal(t, 10*time.Second, DefaultRequestTimeout)
	})
	t.Run("verify_default_flush_interval", func(t *testing.T) {
		assert.Equal(t, 200*time.Millisecond, DefaultFlushInterval)
	})
	t.Run("verify_default_kafka_topic", func(t *testing.T) {
		assert.Equal(t, "de-crypto-events", DefauftKafkaTopic)
	})
	t.Run("verify_default_kafka_brokers", func(t *testing.T) {
		assert.Equal(t, []string{"localhost:9092"}, DefauftKafkaBrokers)
	})
}

func TestGetAddressFile(t *testing.T) {
	t.Run("get_address_file_with_env_variable", func(t *testing.T) {
		os.Setenv("ADDRESS_FILE", "/custom/address.json")
		defer os.Unsetenv("ADDRESS_FILE")
		addressFile := GetAddressFile()
		assert.Equal(t, "/custom/address.json", addressFile)
	})
	t.Run("get_address_file_without_env_variable", func(t *testing.T) {
		os.Unsetenv("ADDRESS_FILE")
		addressFile := GetAddressFile()
		assert.Equal(t, DefaultAddressFile, addressFile)
	})
	t.Run("get_address_file_with_empty_env_variable", func(t *testing.T) {
		os.Setenv("ADDRESS_FILE", "")
		defer os.Unsetenv("ADDRESS_FILE")
		addressFile := GetAddressFile()
		assert.Equal(t, DefaultAddressFile, addressFile)
	})
}

func TestGetKafakTopic(t *testing.T) {
	t.Run("get_kafka_topic_with_env_variable", func(t *testing.T) {
		os.Setenv("KAFKA_TOPIC", "custom-topic")
		defer os.Unsetenv("KAFKA_TOPIC")
		topic := GetKafakTopic()
		assert.Equal(t, "custom-topic", topic)
	})
	t.Run("get_kafka_topic_without_env_variable", func(t *testing.T) {
		os.Unsetenv("KAFKA_TOPIC")
		topic := GetKafakTopic()
		assert.Equal(t, DefauftKafkaTopic, topic)
	})
	t.Run("get_kafka_topic_with_empty_env_variable", func(t *testing.T) {
		os.Setenv("KAFKA_TOPIC", "")
		defer os.Unsetenv("KAFKA_TOPIC")
		topic := GetKafakTopic()
		assert.Equal(t, DefauftKafkaTopic, topic)
	})
}

func TestGetKafakBrokers(t *testing.T) {
	t.Run("get_kafka_brokers_with_env_variable", func(t *testing.T) {
		os.Setenv("KAFKA_BROKERS", "broker1:9092,broker2:9092")
		defer os.Unsetenv("KAFKA_BROKERS")
		brokers := GetKafakBrokers()
		assert.Equal(t, []string{"broker1:9092", "broker2:9092"}, brokers)
	})
	t.Run("get_kafka_brokers_with_single_broker", func(t *testing.T) {
		os.Setenv("KAFKA_BROKERS", "single-broker:9092")
		defer os.Unsetenv("KAFKA_BROKERS")
		brokers := GetKafakBrokers()
		assert.Equal(t, []string{"single-broker:9092"}, brokers)
	})
	t.Run("get_kafka_brokers_without_env_variable", func(t *testing.T) {
		os.Unsetenv("KAFKA_BROKERS")
		brokers := GetKafakBrokers()
		assert.Equal(t, DefauftKafkaBrokers, brokers)
	})
	t.Run("get_kafka_brokers_with_empty_env_variable", func(t *testing.T) {
		os.Setenv("KAFKA_BROKERS", "")
		defer os.Unsetenv("KAFKA_BROKERS")
		brokers := GetKafakBrokers()
		assert.Equal(t, DefauftKafkaBrokers, brokers)
	})
	t.Run("get_kafka_brokers_with_spaces", func(t *testing.T) {
		os.Setenv("KAFKA_BROKERS", " broker1:9092 , broker2:9092 ")
		defer os.Unsetenv("KAFKA_BROKERS")
		brokers := GetKafakBrokers()
		assert.Equal(t, []string{" broker1:9092 ", " broker2:9092 "}, brokers)
	})
}
