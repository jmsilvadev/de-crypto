package internal

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/config"
	"github.com/jmsilvadev/de-crypto/pkg/jsonrpc"
)

type mockAddressIndex struct {
	addresses map[string]string
}

func (m *mockAddressIndex) Lookup(addr string) (string, bool) {
	userID, exists := m.addresses[strings.ToLower(addr)]
	return userID, exists
}

func newMockAddressIndex() *mockAddressIndex {
	return &mockAddressIndex{
		addresses: map[string]string{
			"0xd8da6bf26964af9d7eed9e03e53415d37aa96045": "vitalik",
			"0x742d35cc6634c0532925a3b8d0c0c2c0c0c0c0c0": "binance",
			"0x28c6c06298d514db089934071355e5743bf21d60": "binance2",
		},
	}
}

func TestFilterMatcher(t *testing.T) {
	t.Run("filter_matcher_processes_blocks", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		blocksCh := make(chan jsonrpc.Block, 1)
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		cfg := config.FilterConfig{
			Workers: 1,
		}
		go filterMatcher(ctx, cfg, blocksCh, addrIdx, eventsCh)
		toAddr := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
		block := jsonrpc.Block{
			Number: "0x3039",
			Hash:   "0xabc123",
			Transactions: []jsonrpc.Transaction{
				{
					From:  "0x1234567890123456789012345678901234567890",
					To:    &toAddr,
					Value: "0xde0b6b3a7640000",
					Hash:  "0xtx123",
				},
			},
		}
		blocksCh <- block
		close(blocksCh)
		select {
		case event := <-eventsCh:
			if event.UserID != "vitalik" {
				t.Errorf("Expected userID 'vitalik', got '%s'", event.UserID)
			}
		case <-time.After(30 * time.Millisecond):
			t.Error("Expected event but got timeout")
		}
	})
	t.Run("filter_matcher_handles_context_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		blocksCh := make(chan jsonrpc.Block, 1)
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		cfg := config.FilterConfig{
			Workers: 1,
		}
		go filterMatcher(ctx, cfg, blocksCh, addrIdx, eventsCh)
		cancel()
	})
	t.Run("filter_matcher_with_zero_workers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		blocksCh := make(chan jsonrpc.Block, 1)
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		cfg := config.FilterConfig{
			Workers: 0,
		}
		go filterMatcher(ctx, cfg, blocksCh, addrIdx, eventsCh)
	})
	t.Run("filter_matcher_with_negative_workers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		blocksCh := make(chan jsonrpc.Block, 1)
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		cfg := config.FilterConfig{
			Workers: -1,
		}
		go filterMatcher(ctx, cfg, blocksCh, addrIdx, eventsCh)
	})
	t.Run("filter_matcher_channel_close", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		blocksCh := make(chan jsonrpc.Block, 1)
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		cfg := config.FilterConfig{
			Workers: 1,
		}
		go filterMatcher(ctx, cfg, blocksCh, addrIdx, eventsCh)
		close(blocksCh)
	})
}

func TestProcessBlock(t *testing.T) {
	t.Run("process_block_with_matching_from_address", func(t *testing.T) {
		ctx := context.Background()
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		toAddr := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		block := jsonrpc.Block{
			Number: "0x3039",
			Hash:   "0xabc123",
			Transactions: []jsonrpc.Transaction{
				{
					From:  "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
					To:    &toAddr,
					Value: "0xde0b6b3a7640000",
					Hash:  "0xtx123",
				},
			},
		}
		processBlock(ctx, block, addrIdx, eventsCh)
		select {
		case event := <-eventsCh:
			if event.UserID != "vitalik" {
				t.Errorf("Expected userID 'vitalik', got '%s'", event.UserID)
			}
		case <-time.After(10 * time.Millisecond):
			t.Error("Expected event but got timeout")
		}
		close(eventsCh)
	})
	t.Run("process_block_with_matching_to_address", func(t *testing.T) {
		ctx := context.Background()
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		toAddr := "0x742d35Cc6634C0532925a3b8D0C0C2C0C0C0C0C0"
		block := jsonrpc.Block{
			Number: "0x3039",
			Hash:   "0xabc123",
			Transactions: []jsonrpc.Transaction{
				{
					From:  "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
					To:    &toAddr,
					Value: "0xde0b6b3a7640000",
					Hash:  "0xtx123",
				},
			},
		}
		processBlock(ctx, block, addrIdx, eventsCh)
		select {
		case event := <-eventsCh:
			if event.UserID != "binance" {
				t.Errorf("Expected userID 'binance', got '%s'", event.UserID)
			}
		case <-time.After(10 * time.Millisecond):
			t.Error("Expected event but got timeout")
		}
		close(eventsCh)
	})
	t.Run("process_block_with_no_matching_addresses", func(t *testing.T) {
		ctx := context.Background()
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		toAddr := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
		block := jsonrpc.Block{
			Number: "0x3039",
			Hash:   "0xabc123",
			Transactions: []jsonrpc.Transaction{
				{
					From:  "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
					To:    &toAddr,
					Value: "0xde0b6b3a7640000",
					Hash:  "0xtx123",
				},
			},
		}
		processBlock(ctx, block, addrIdx, eventsCh)
		select {
		case <-eventsCh:
			t.Error("Expected no events but got one")
		case <-time.After(10 * time.Millisecond):
		}
		close(eventsCh)
	})
	t.Run("process_block_with_nil_to_address", func(t *testing.T) {
		ctx := context.Background()
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		block := jsonrpc.Block{
			Number: "0x3039",
			Hash:   "0xabc123",
			Transactions: []jsonrpc.Transaction{
				{
					From:  "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
					To:    nil,
					Value: "0xde0b6b3a7640000",
					Hash:  "0xtx123",
				},
			},
		}
		processBlock(ctx, block, addrIdx, eventsCh)
		select {
		case event := <-eventsCh:
			if event.UserID != "vitalik" {
				t.Errorf("Expected userID 'vitalik', got '%s'", event.UserID)
			}
		case <-time.After(10 * time.Millisecond):
			t.Error("Expected event but got timeout")
		}
		close(eventsCh)
	})
	t.Run("process_block_with_invalid_block_number", func(t *testing.T) {
		ctx := context.Background()
		eventsCh := make(chan Event, 1)
		addrIdx := newMockAddressIndex()
		toAddr := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
		block := jsonrpc.Block{
			Number: "invalid",
			Hash:   "0xabc123",
			Transactions: []jsonrpc.Transaction{
				{
					From:  "0x1234567890123456789012345678901234567890",
					To:    &toAddr,
					Value: "0xde0b6b3a7640000",
					Hash:  "0xtx123",
				},
			},
		}
		processBlock(ctx, block, addrIdx, eventsCh)
		select {
		case <-eventsCh:
			t.Error("Expected no events due to invalid block number but got one")
		case <-time.After(10 * time.Millisecond):
		}
		close(eventsCh)
	})
}

func TestEmitEvent(t *testing.T) {
	t.Run("emit_event_success", func(t *testing.T) {
		ctx := context.Background()
		eventsCh := make(chan Event, 1)
		event := Event{
			BlockNumber: 12345,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		emitEvent(ctx, eventsCh, event)
		select {
		case receivedEvent := <-eventsCh:
			if receivedEvent.UserID != "user1" {
				t.Errorf("Expected userID 'user1', got '%s'", receivedEvent.UserID)
			}
		case <-time.After(10 * time.Millisecond):
			t.Error("Expected event but got timeout")
		}
		close(eventsCh)
	})
	t.Run("emit_event_context_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		eventsCh := make(chan Event, 1)
		event := Event{
			BlockNumber: 12345,
			TxHash:      "0xabc123",
			From:        "0x1234567890123456789012345678901234567890",
			To:          "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			UserID:      "user1",
		}
		cancel()
		emitEvent(ctx, eventsCh, event)
		select {
		case <-eventsCh:
		case <-time.After(10 * time.Millisecond):
		}
		close(eventsCh)
	})
}
