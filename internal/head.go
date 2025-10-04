package internal

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/config"
	"github.com/jmsilvadev/de-crypto/pkg/jsonrpc"
)

func headMonitor(ctx context.Context, cfg config.HeadMonitorConfig, rpc jsonrpc.JsonRpcClient, headsCh chan<- uint64) {
	log.Println("Starting head monitor")

	nextHeight := cfg.StartFrom
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			head, err := rpc.GetCurrentBlockNumber(ctx)
			if err == nil {
				sent := 0
				for nextHeight <= head && sent < cfg.MaxEnqueuePerTick {
					select {
					case <-ctx.Done():
						return
					case headsCh <- nextHeight:
						nextHeight++
						sent++
					default:
						sent = cfg.MaxEnqueuePerTick
					}
				}
			}

			// lets add a jit to avoid have the amount of request at the same time
			// we can also reduce this time if the queue is full but i'll maintain the code simpler
			delay := withJitter(cfg.PollInterval, cfg.Jitter)
			resetTimer(timer, delay)
		}
	}
}

func withJitter(base time.Duration, jitter float64) time.Duration {
	if jitter <= 0 {
		return base
	}
	fmin, fmax := 1.0-jitter, 1.0+jitter
	f := fmin + rand.Float64()*(fmax-fmin)
	return time.Duration(float64(base) * f)
}

func resetTimer(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(d)
}
