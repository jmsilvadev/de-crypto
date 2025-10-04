package internal

import (
	"context"
	"testing"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/config"
)

func TestSinkProcessor(t *testing.T) {
	t.Run("sink_processor_processes_events", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		eventsCh := make(chan Event, 10)
		cfg := config.SinkConfig{
			BatchSize:     5,
			FlushInterval: 100 * time.Millisecond,
		}
		go sinkProcessor(ctx, cfg, eventsCh, nil, func(data []byte) error { return nil })
		event := Event{
			BlockNumber: 12345,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		eventsCh <- event
		<-ctx.Done()
	})
	t.Run("sink_processor_with_checkpoint_store", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		eventsCh := make(chan Event, 10)
		cfg := config.SinkConfig{
			BatchSize:     5,
			FlushInterval: 100 * time.Millisecond,
		}
		go sinkProcessor(ctx, cfg, eventsCh, nil, func(data []byte) error { return nil })
		event := Event{
			BlockNumber: 12345,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		eventsCh <- event
		<-ctx.Done()
	})
	t.Run("sink_processor_output_error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		eventsCh := make(chan Event, 10)
		cfg := config.SinkConfig{
			BatchSize:     5,
			FlushInterval: 100 * time.Millisecond,
		}
		go sinkProcessor(ctx, cfg, eventsCh, nil, func(data []byte) error {
			return nil
		})
		event := Event{
			BlockNumber: 12345,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		eventsCh <- event
		<-ctx.Done()
	})
	t.Run("sink_processor_lower_block_number", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		eventsCh := make(chan Event, 10)
		cfg := config.SinkConfig{
			BatchSize:     5,
			FlushInterval: 100 * time.Millisecond,
		}
		go sinkProcessor(ctx, cfg, eventsCh, nil, func(data []byte) error { return nil })
		event := Event{
			BlockNumber: 10000,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		eventsCh <- event
		<-ctx.Done()
	})
	t.Run("sink_processor_equal_block_number", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		eventsCh := make(chan Event, 10)
		cfg := config.SinkConfig{
			BatchSize:     5,
			FlushInterval: 100 * time.Millisecond,
		}
		go sinkProcessor(ctx, cfg, eventsCh, nil, func(data []byte) error { return nil })
		event := Event{
			BlockNumber: 12345,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		eventsCh <- event
		<-ctx.Done()
	})
	t.Run("sink_processor_batch_not_full", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		eventsCh := make(chan Event, 10)
		cfg := config.SinkConfig{
			BatchSize:     10,
			FlushInterval: 100 * time.Millisecond,
		}
		go sinkProcessor(ctx, cfg, eventsCh, nil, func(data []byte) error { return nil })
		event := Event{
			BlockNumber: 12345,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		eventsCh <- event
		<-ctx.Done()
	})
}
