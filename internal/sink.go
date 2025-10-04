package internal

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/checkpoint"
	"github.com/jmsilvadev/de-crypto/pkg/config"
)

func sinkProcessor(ctx context.Context, cfg config.SinkConfig, eventsCh <-chan Event, store *checkpoint.CheckpointStore, publisher func([]byte) error) {
	log.Println("Starting sink processor")

	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = 250 * time.Millisecond
	}

	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 256
	}

	flushTicker := time.NewTicker(cfg.FlushInterval)
	defer flushTicker.Stop()

	checkpointTicker := time.NewTicker(500 * time.Millisecond)
	defer checkpointTicker.Stop()

	var lastSaved uint64
	var pendingCheckpoint *uint64
	var maxSeen uint64

	if store != nil {
		if n, err := store.Load(); err != nil {
			log.Println(err)
		} else {
			lastSaved = n
			maxSeen = n
		}
	}

	// using literals here to reuse the state...
	batch := make([][]byte, 0, cfg.BatchSize)

	sendEvents := func() {
		if len(batch) == 0 {
			return
		}
		for _, msg := range batch {
			if err := publisher(msg); err != nil {
				// so if we have an error we will continue and log
				// we can create an improvement to add a dead letter or so...
				log.Println("publish:", err)
			}
		}
		batch = batch[:0]
	}

	// same here
	saveCheckpointIfNeeded := func() {
		if store == nil || pendingCheckpoint == nil {
			return
		}
		confirmed := *pendingCheckpoint
		if confirmed > lastSaved {
			if err := store.Save(confirmed); err != nil {
				log.Println(err)
			} else {
				lastSaved = confirmed
			}
		}
		pendingCheckpoint = nil
	}

	//same hre
	writeEvent := func(ev Event) {
		b, err := json.Marshal(ev)
		if err != nil {
			log.Println(err)
			return
		}
		batch = append(batch, b)

		if ev.BlockNumber > maxSeen {
			maxSeen = ev.BlockNumber
			if maxSeen > lastSaved {
				pc := maxSeen
				pendingCheckpoint = &pc
			}
		}

		if len(batch) >= cfg.BatchSize {
			sendEvents()
		}
	}

	for {
		select {
		case <-ctx.Done():
			sendEvents()
			saveCheckpointIfNeeded()
			return
		case ev, ok := <-eventsCh:
			if !ok {
				sendEvents()
				saveCheckpointIfNeeded()
				return
			}
			writeEvent(ev)
		case <-flushTicker.C:
			sendEvents()
		case <-checkpointTicker.C:
			saveCheckpointIfNeeded()
		}
	}
}
