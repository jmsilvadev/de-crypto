package internal

import (
	"context"
	"errors"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/config"
	"github.com/jmsilvadev/de-crypto/pkg/jsonrpc"
	"github.com/jmsilvadev/de-crypto/pkg/utils"
)

func blockFetcher(ctx context.Context, cfg config.BlockFetcherConfig, rpc jsonrpc.JsonRpcClient, headCh <-chan uint64, out chan<- jsonrpc.Block) {
	for i := 0; i < cfg.Workers; i++ {
		go worker(ctx, rpc, cfg, headCh, out)
	}
	<-ctx.Done()
}

func worker(ctx context.Context, rpc jsonrpc.JsonRpcClient, cfg config.BlockFetcherConfig, headCh <-chan uint64, out chan<- jsonrpc.Block) {
	log.Println("Starting block fetcher worker")
	for {
		select {
		case <-ctx.Done():
			return
		case h, ok := <-headCh:
			if !ok {
				return
			}
			blk, err := fetchWithRetry(ctx, rpc, cfg, h)
			if err != nil {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- *blk:
			}
		}
	}
}

func fetchWithRetry(ctx context.Context, rpc jsonrpc.JsonRpcClient, cfg config.BlockFetcherConfig, blockNumber uint64) (*jsonrpc.Block, error) {
	attempt := 0
	for {
		reqCtx, cancel := context.WithTimeout(ctx, cfg.ReqTimeout)
		blk, err := rpc.GetBlockByNumber(reqCtx, blockNumber)
		cancel()

		if err != nil {
			return nil, err
		}

		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ctx.Err()
		}

		blkNumber, err := utils.ParseHexUint64(blk.Number)
		if err != nil {
			return nil, err
		}

		if blk != nil && blkNumber == blockNumber {
			return blk, nil
		}

		delay := backoffWithJitter(cfg, attempt)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
		attempt++
	}
}

func backoffWithJitter(cfg config.BlockFetcherConfig, attempt int) time.Duration {
	exp := float64(cfg.RetryBaseDelay) * math.Pow(2, float64(attempt)) // veru basic backoff
	base := time.Duration(exp)
	if base > cfg.RetryMaxDelay {
		base = cfg.RetryMaxDelay
	}
	if cfg.Jitter <= 0 {
		return base
	}

	fmin, fmax := 1.0-cfg.Jitter, 1.0+cfg.Jitter
	f := fmin + rand.Float64()*(fmax-fmin)
	return time.Duration(float64(base) * f)
}
