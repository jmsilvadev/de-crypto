package internal

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jmsilvadev/de-crypto/pkg/address"
	"github.com/jmsilvadev/de-crypto/pkg/checkpoint"
	"github.com/jmsilvadev/de-crypto/pkg/config"
	"github.com/jmsilvadev/de-crypto/pkg/jsonrpc"
	"github.com/jmsilvadev/de-crypto/pkg/kafka"
)

func Start() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	store := checkpoint.NewCheckpointStore(config.DefaultCheckpointStore)

	confirmedCheckpointFromDisk, err := store.Load()
	if err != nil {
		// lets assume start with zero to not have to stop the system and only log the reason here
		log.Println(err)
	}

	// To understand under 15 minutes
	headsCh := make(chan uint64, config.DefaultHeadsChannelSize)          // chanell to get and buffer the current blocks ans send to the blocks fetcher
	blocksCh := make(chan jsonrpc.Block, config.DefaultBlocksChannelSize) // channel do get the full blocks details and send to the filter matcher
	eventsCh := make(chan Event, config.DefaultEventsChannelSize)         // channel to send the filtered blocks to sink

	cfgHead := config.HeadMonitorConfig{
		PollInterval:      config.DefaultPollingInterval,
		StartFrom:         confirmedCheckpointFromDisk,
		Jitter:            config.DeaultJitterDeviation,
		MaxEnqueuePerTick: config.DefaultMaxEnqueuePerTick,
	}

	jsonRPC := jsonrpc.NewEthereum(config.CliUrl, config.DefaultHttpClient)

	wg.Add(1)
	go func() {
		defer wg.Done()
		headMonitor(ctx, cfgHead, jsonRPC, headsCh)
	}()

	cfgBlock := config.BlockFetcherConfig{
		Workers:        config.DefaultWorkersQty,
		ReqTimeout:     config.DefaultRequestTimeout,
		RetryBaseDelay: config.DefaultPollingInterval,
		RetryMaxDelay:  config.DefaultPollingInterval,
		Jitter:         config.DeaultJitterDeviation,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		blockFetcher(ctx, cfgBlock, jsonRPC, headsCh, blocksCh)
	}()

	addressFile := config.GetAddressFile()
	add, err := address.NewMemoryAddressIndexFromJSON(addressFile)
	if err != nil {
		panic(err)
	}

	cfgFilter := config.FilterConfig{
		Workers: 8,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		filterMatcher(ctx, cfgFilter, blocksCh, add, eventsCh)
	}()

	cfgSink := config.SinkConfig{
		FlushInterval: config.DefaultFlushInterval,
		BatchSize:     config.DefaultBatchSize,
	}

	cfgKafka := config.KafkaConfig{
		Brokers: config.GetKafakBrokers(),
		Topic:   config.GetKafakTopic(),
	}

	pub, err := kafka.NewPublisher(ctx, cfgKafka.Brokers, cfgKafka.Topic)
	if err != nil {
		panic(err)
	}
	defer pub.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sinkProcessor(ctx, cfgSink, eventsCh, store, pub.Publish)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	cancel()

	<-sigCh
	os.Exit(1)

	wg.Wait()
}
