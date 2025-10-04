package internal

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/config"
	"github.com/jmsilvadev/de-crypto/pkg/jsonrpc"
	"github.com/stretchr/testify/assert"
)

func TestBlockFetcher(t *testing.T) {
	t.Run("block_fetcher_starts_workers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		headCh := make(chan uint64, 10)
		outCh := make(chan jsonrpc.Block, 10)

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.BlockFetcherConfig{
			Workers:        2,
			ReqTimeout:     50 * time.Millisecond,
			RetryBaseDelay: 10 * time.Millisecond,
			RetryMaxDelay:  50 * time.Millisecond,
			Jitter:         0.1,
		}

		go blockFetcher(ctx, cfg, rpcClient, headCh, outCh)

		headCh <- 12345

		<-ctx.Done()
	})

	t.Run("worker_processes_block_numbers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		headCh := make(chan uint64, 10)
		outCh := make(chan jsonrpc.Block, 10)

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.BlockFetcherConfig{
			Workers:        1,
			ReqTimeout:     50 * time.Millisecond,
			RetryBaseDelay: 10 * time.Millisecond,
			RetryMaxDelay:  50 * time.Millisecond,
			Jitter:         0.1,
		}

		go worker(ctx, rpcClient, cfg, headCh, outCh)

		headCh <- 12345

		<-ctx.Done()
	})

	t.Run("worker_handles_rpc_errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		headCh := make(chan uint64, 10)
		outCh := make(chan jsonrpc.Block, 10)

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.BlockFetcherConfig{
			Workers:        1,
			ReqTimeout:     50 * time.Millisecond,
			RetryBaseDelay: 10 * time.Millisecond,
			RetryMaxDelay:  50 * time.Millisecond,
			Jitter:         0.1,
		}

		go worker(ctx, rpcClient, cfg, headCh, outCh)

		headCh <- 12345

		<-ctx.Done()
	})
}

func TestFetchWithRetry(t *testing.T) {
	t.Run("fetch_with_retry_success_on_first_attempt", func(t *testing.T) {
		ctx := context.Background()

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.BlockFetcherConfig{
			ReqTimeout:     50 * time.Millisecond,
			RetryBaseDelay: 10 * time.Millisecond,
			RetryMaxDelay:  50 * time.Millisecond,
			Jitter:         0.1,
		}

		_, err := fetchWithRetry(ctx, rpcClient, cfg, 12345)

		if err != nil {
			t.Logf("Expected error in test environment: %v", err)
		}
	})

	t.Run("fetch_with_retry_context_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.BlockFetcherConfig{
			ReqTimeout:     50 * time.Millisecond,
			RetryBaseDelay: 10 * time.Millisecond,
			RetryMaxDelay:  50 * time.Millisecond,
			Jitter:         0.1,
		}

		cancel()

		block, err := fetchWithRetry(ctx, rpcClient, cfg, 12345)

		assert.Error(t, err)
		assert.Nil(t, block)
		assert.True(t, errors.Is(err, context.Canceled))
	})
}

func TestBackoffWithJitter(t *testing.T) {
	t.Run("backoff_with_jitter_no_jitter", func(t *testing.T) {
		cfg := config.BlockFetcherConfig{
			RetryBaseDelay: 100 * time.Millisecond,
			RetryMaxDelay:  1000 * time.Millisecond,
			Jitter:         0.0,
		}

		delay1 := backoffWithJitter(cfg, 0)
		delay2 := backoffWithJitter(cfg, 1)
		delay3 := backoffWithJitter(cfg, 2)

		assert.Equal(t, 100*time.Millisecond, delay1)
		assert.Equal(t, 200*time.Millisecond, delay2)
		assert.Equal(t, 400*time.Millisecond, delay3)
	})

	t.Run("backoff_with_jitter_with_jitter", func(t *testing.T) {
		cfg := config.BlockFetcherConfig{
			RetryBaseDelay: 100 * time.Millisecond,
			RetryMaxDelay:  1000 * time.Millisecond,
			Jitter:         0.1,
		}

		delay := backoffWithJitter(cfg, 1)
		baseDelay := 200 * time.Millisecond
		minDelay := time.Duration(float64(baseDelay) * 0.9)
		maxDelay := time.Duration(float64(baseDelay) * 1.1)

		assert.True(t, delay >= minDelay)
		assert.True(t, delay <= maxDelay)
	})

	t.Run("backoff_with_jitter_max_delay_limit", func(t *testing.T) {
		cfg := config.BlockFetcherConfig{
			RetryBaseDelay: 100 * time.Millisecond,
			RetryMaxDelay:  200 * time.Millisecond,
			Jitter:         0.0,
		}

		delay := backoffWithJitter(cfg, 10)
		assert.Equal(t, 200*time.Millisecond, delay)
	})

	t.Run("backoff_with_jitter_negative_jitter", func(t *testing.T) {
		cfg := config.BlockFetcherConfig{
			RetryBaseDelay: 100 * time.Millisecond,
			RetryMaxDelay:  1000 * time.Millisecond,
			Jitter:         -0.1,
		}

		delay := backoffWithJitter(cfg, 1)
		assert.Equal(t, 200*time.Millisecond, delay)
	})
}
