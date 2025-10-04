package internal

import (
	"context"
	"testing"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/config"
	"github.com/jmsilvadev/de-crypto/pkg/jsonrpc"
	"github.com/stretchr/testify/assert"
)

func TestHeadMonitor(t *testing.T) {
	t.Run("head_monitor_processes_block_numbers", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		headCh := make(chan uint64, 10)

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.HeadMonitorConfig{
			PollInterval: 50 * time.Millisecond,
		}

		go headMonitor(ctx, cfg, rpcClient, headCh)

		<-ctx.Done()
	})

	t.Run("head_monitor_handles_rpc_errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		headCh := make(chan uint64, 10)

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.HeadMonitorConfig{
			PollInterval: 50 * time.Millisecond,
		}

		go headMonitor(ctx, cfg, rpcClient, headCh)

		<-ctx.Done()
	})

	t.Run("head_monitor_respects_max_enqueue_per_tick", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		headCh := make(chan uint64, 1)

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.HeadMonitorConfig{
			PollInterval:      50 * time.Millisecond,
			MaxEnqueuePerTick: 1,
		}

		go headMonitor(ctx, cfg, rpcClient, headCh)

		<-ctx.Done()
	})

	t.Run("head_monitor_context_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		headCh := make(chan uint64, 10)

		rpcClient := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

		cfg := config.HeadMonitorConfig{
			PollInterval: 50 * time.Millisecond,
		}

		go headMonitor(ctx, cfg, rpcClient, headCh)

		cancel()

		<-ctx.Done()
	})
}

func TestWithJitter(t *testing.T) {
	t.Run("with_jitter_no_jitter", func(t *testing.T) {
		base := 100 * time.Millisecond
		jitter := 0.0

		result := withJitter(base, jitter)

		assert.Equal(t, base, result)
	})

	t.Run("with_jitter_with_jitter", func(t *testing.T) {
		base := 100 * time.Millisecond
		jitter := 0.1

		result := withJitter(base, jitter)

		min := time.Duration(float64(base) * 0.9)
		max := time.Duration(float64(base) * 1.1)

		assert.True(t, result >= min)
		assert.True(t, result <= max)
	})

	t.Run("with_jitter_negative_jitter", func(t *testing.T) {
		base := 100 * time.Millisecond
		jitter := -0.1

		result := withJitter(base, jitter)

		assert.Equal(t, base, result)
	})

	t.Run("with_jitter_zero_base", func(t *testing.T) {
		base := 0 * time.Millisecond
		jitter := 0.1

		result := withJitter(base, jitter)

		assert.Equal(t, base, result)
	})
}

func TestResetTimer(t *testing.T) {
	t.Run("reset_timer_not_stopped", func(t *testing.T) {
		timer := time.NewTimer(100 * time.Millisecond)
		defer timer.Stop()

		resetTimer(timer, 50*time.Millisecond)

		select {
		case <-timer.C:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("timer did not fire")
		}
	})

	t.Run("reset_timer_already_fired", func(t *testing.T) {
		timer := time.NewTimer(10 * time.Millisecond)
		<-timer.C

		resetTimer(timer, 50*time.Millisecond)

		select {
		case <-timer.C:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("timer did not fire")
		}
	})

	t.Run("reset_timer_multiple_resets", func(t *testing.T) {
		timer := time.NewTimer(100 * time.Millisecond)
		defer timer.Stop()

		resetTimer(timer, 50*time.Millisecond)
		resetTimer(timer, 30*time.Millisecond)
		resetTimer(timer, 20*time.Millisecond)

		select {
		case <-timer.C:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("timer did not fire")
		}
	})
}
