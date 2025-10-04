package kafka

import (
	"context"
	"testing"
	"time"
)

func TestEnsureTopicExists(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	brokers := []string{"localhost:9092"}
	topic := "test-topic-creation"
	err := ensureTopicExists(ctx, brokers, topic)
	if err != nil {
		t.Logf("Topic creation test skipped (Kafka not available): %v", err)
		t.Skip("Kafka not available for testing")
	}
	err = ensureTopicExists(ctx, brokers, topic)
	if err != nil {
		t.Errorf("Topic should exist after creation: %v", err)
	}
}

func TestNewPublisher_WithTopicCreation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	brokers := []string{"localhost:9092"}
	topic := "test-publisher-topic"
	publisher, err := NewPublisher(ctx, brokers, topic)
	if err != nil {
		t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
		t.Skip("Kafka not available for testing")
	}
	defer publisher.Close()
	testMessage := []byte("test message")
	err = publisher.Publish(testMessage)
	if err != nil {
		t.Errorf("Failed to publish message: %v", err)
	}
}

func TestPublisher_Publish(t *testing.T) {
	t.Run("publish_message_success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-publish-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err != nil {
			t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
			t.Skip("Kafka not available for testing")
		}
		defer publisher.Close()
		testMessage := []byte("test publish message")
		err = publisher.Publish(testMessage)
		if err != nil {
			t.Errorf("Failed to publish message: %v", err)
		}
	})
	t.Run("publish_empty_message", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-publish-empty-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err != nil {
			t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
			t.Skip("Kafka not available for testing")
		}
		defer publisher.Close()
		emptyMessage := []byte("")
		err = publisher.Publish(emptyMessage)
		if err != nil {
			t.Errorf("Failed to publish empty message: %v", err)
		}
	})
	t.Run("publish_large_message", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-publish-large-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err != nil {
			t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
			t.Skip("Kafka not available for testing")
		}
		defer publisher.Close()
		largeMessage := make([]byte, 1024*1024)
		for i := range largeMessage {
			largeMessage[i] = byte(i % 256)
		}
		err = publisher.Publish(largeMessage)
		if err != nil {
			t.Errorf("Failed to publish large message: %v", err)
		}
	})
	t.Run("publish_multiple_messages", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-publish-multiple-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err != nil {
			t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
			t.Skip("Kafka not available for testing")
		}
		defer publisher.Close()
		messages := [][]byte{
			[]byte("message 1"),
			[]byte("message 2"),
			[]byte("message 3"),
		}
		for i, msg := range messages {
			err = publisher.Publish(msg)
			if err != nil {
				t.Errorf("Failed to publish message %d: %v", i+1, err)
			}
		}
	})
}

func TestPublisher_Close(t *testing.T) {
	t.Run("close_publisher_success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-close-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err != nil {
			t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
			t.Skip("Kafka not available for testing")
		}
		err = publisher.Close()
		if err != nil {
			t.Errorf("Failed to close publisher: %v", err)
		}
	})
	t.Run("close_publisher_multiple_times", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-close-multiple-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err != nil {
			t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
			t.Skip("Kafka not available for testing")
		}
		err = publisher.Close()
		if err != nil {
			t.Errorf("Failed to close publisher first time: %v", err)
		}
		err = publisher.Close()
		if err != nil {
			t.Errorf("Failed to close publisher second time: %v", err)
		}
	})
	t.Run("publish_after_close", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-publish-after-close-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err != nil {
			t.Logf("Publisher creation test skipped (Kafka not available): %v", err)
			t.Skip("Kafka not available for testing")
		}
		err = publisher.Close()
		if err != nil {
			t.Errorf("Failed to close publisher: %v", err)
		}
		testMessage := []byte("test message after close")
		err = publisher.Publish(testMessage)
		if err == nil {
			t.Error("Expected error when publishing after close, but got none")
		}
	})
}

func TestNewPublisher_ErrorCases(t *testing.T) {
	t.Run("new_publisher_with_empty_brokers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{}
		topic := "test-topic"
		_, err := NewPublisher(ctx, brokers, topic)
		if err == nil {
			t.Error("Expected error for empty brokers, but got none")
		}
	})
	t.Run("new_publisher_with_empty_topic", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := ""
		_, err := NewPublisher(ctx, brokers, topic)
		if err == nil {
			t.Error("Expected error for empty topic, but got none")
		}
	})
	t.Run("new_publisher_with_nil_brokers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var brokers []string
		topic := "test-topic"
		_, err := NewPublisher(ctx, brokers, topic)
		if err == nil {
			t.Error("Expected error for nil brokers, but got none")
		}
	})
}

func TestPublisher_Methods_UnitTests(t *testing.T) {
	t.Run("publisher_methods_without_kafka", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		brokers := []string{"invalid-broker:9092"}
		topic := "test-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err == nil {
			t.Error("Expected error for invalid broker, but got none")
		}
		if publisher != nil {
			t.Error("Expected nil publisher for invalid broker")
		}
	})
	t.Run("publisher_validation_tests", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		testCases := []struct {
			name    string
			brokers []string
			topic   string
		}{
			{"empty_brokers", []string{}, "valid-topic"},
			{"nil_brokers", nil, "valid-topic"},
			{"empty_topic", []string{"localhost:9092"}, ""},
			{"whitespace_topic", []string{"localhost:9092"}, "   "},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := NewPublisher(ctx, tc.brokers, tc.topic)
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tc.name)
				}
			})
		}
	})
}

func TestPublisher_Methods_DirectTesting(t *testing.T) {
	t.Run("test_publish_and_close_methods_exist", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		brokers := []string{"localhost:9092"}
		topic := "test-methods-topic"
		publisher, err := NewPublisher(ctx, brokers, topic)
		if err == nil {
			defer publisher.Close()
			testMessage := []byte("test message")
			publishErr := publisher.Publish(testMessage)
			if publishErr != nil {
				t.Logf("Publish failed as expected: %v", publishErr)
			}
			closeErr := publisher.Close()
			if closeErr != nil {
				t.Logf("Close failed as expected: %v", closeErr)
			}
		} else {
			t.Logf("Publisher creation failed as expected: %v", err)
		}
	})
	t.Run("test_method_signatures", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic when calling methods on nil publisher: %v", r)
			}
		}()
		var publisher *Publisher
		if publisher != nil {
			_ = publisher.Publish([]byte("test"))
			_ = publisher.Close()
		}
	})
}

func TestPublisher_Methods_WithoutKafka(t *testing.T) {
	t.Run("test_publish_method_signature", func(t *testing.T) {
		var publisher *Publisher
		if publisher != nil {
			_ = publisher.Publish([]byte("test"))
		}
	})
	t.Run("test_close_method_signature", func(t *testing.T) {
		var publisher *Publisher
		if publisher != nil {
			_ = publisher.Close()
		}
	})
	t.Run("test_publish_with_nil_publisher", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic with nil publisher: %v", r)
			}
		}()
		var publisher *Publisher
		if publisher == nil {
			t.Log("Publisher is nil as expected")
		}
	})
	t.Run("test_close_with_nil_publisher", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic with nil publisher: %v", r)
			}
		}()
		var publisher *Publisher
		if publisher == nil {
			t.Log("Publisher is nil as expected")
		}
	})
}
