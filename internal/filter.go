package internal

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/jmsilvadev/de-crypto/pkg/address"
	"github.com/jmsilvadev/de-crypto/pkg/config"
	"github.com/jmsilvadev/de-crypto/pkg/jsonrpc"
	"github.com/jmsilvadev/de-crypto/pkg/utils"
)

func filterMatcher(ctx context.Context, cfg config.FilterConfig, blocksCh <-chan jsonrpc.Block, addrIdx address.AddressIndex, eventsCh chan<- Event) {
	log.Println("Starting filter matcher")

	workers := cfg.Workers
	if workers <= 0 {
		workers = 1
	}

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case b, ok := <-blocksCh:
					if !ok {
						return
					}
					processBlock(ctx, b, addrIdx, eventsCh)
				}
			}
		}()
	}

	wg.Wait()
}

func processBlock(ctx context.Context, b jsonrpc.Block, addrIdx address.AddressIndex, eventsCh chan<- Event) {
	for _, tx := range b.Transactions {
		from := strings.ToLower(tx.From)
		to := ""
		if tx.To != nil {
			to = strings.ToLower(*tx.To)
		}

		n, err := utils.ParseHexUint64(b.Number)
		if err != nil {
			log.Println(err)
			continue
		}

		if userID, ok := addrIdx.Lookup(from); ok {
			ev := Event{
				UserID:      userID,
				From:        tx.From,
				To:          to,
				AmountWei:   tx.Value,
				TxHash:      tx.Hash,
				BlockNumber: n,
			}
			emitEvent(ctx, eventsCh, ev)
		}

		if userID, ok := addrIdx.Lookup(to); ok {
			ev := Event{
				UserID:      userID,
				From:        tx.From,
				To:          to,
				AmountWei:   tx.Value,
				TxHash:      tx.Hash,
				BlockNumber: n,
			}
			emitEvent(ctx, eventsCh, ev)
		}
	}
}

func emitEvent(ctx context.Context, eventsCh chan<- Event, ev Event) {
	select {
	case <-ctx.Done():
		return
	case eventsCh <- ev:
	}
}
