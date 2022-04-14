package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"

	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming/coinbase"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming/coinbase/handler"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultPort is the default websocket URL to subscribe to.
	DefaultWebSocketURL = "wss://ws-feed.exchange.coinbase.com"
	// DefaultLogLevel is the default log level set for logrus.
	DefaultLogLevel = logrus.FatalLevel
	// DefaultPairs is the default list of pairs to get the vwap for.
	DefaultPairs = "BTC-USD,ETH-USD,ETH-BTC"
	// DefaultVwapWindowSize is the default window size for the vwap calculation.
	DefaultVwapWindowSize = 200
)

func main() {
	var (
		queryPairs     = flag.String("pairs", DefaultPairs, "comma separated list of pairs to query")
		verbose        = flag.Bool("verbose", false, "verbose logging")
		wsURL          = flag.String("wsurl", DefaultWebSocketURL, "websocket url")
		vwapWindowSize = flag.Int("window-size", DefaultVwapWindowSize, "vwap window size")
	)

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Set up logger.
	logger := logrus.New()
	logger.SetLevel(DefaultLogLevel)
	if *verbose {
		logger.SetLevel(logrus.TraceLevel)
	}

	productIds := strings.Split(*queryPairs, ",")

	// Build the request to subscribe to the coinbase websocket feed.
	subscribeReq := coinbase.SubscribeRequest{
		Type:       "subscribe",
		ProductIds: productIds,
		Channels: []coinbase.Channel{
			{
				Name:       "matches",
				ProductIds: productIds,
			},
		},
	}

	request, err := json.Marshal(subscribeReq)
	if err != nil {
		logger.Errorf("failed to marshal subscribe request: %v", err)
	}

	ctx := context.Background()

	// Create a new websocket streaming client.
	streamer := coinbase.NewStreamer(ctx, *wsURL, string(request))
	streamer.SetLogger(logger)

	var streamHandler streaming.StreamDataHandler

	// Create a new vwap data handler.
	streamHandler = handler.NewStreamDataHandler(*vwapWindowSize, productIds)
	streamHandler.SetLogger(logger)
	streamHandler.SetStreamer(streamer)

	logger.Infoln("Starting vwap price streaming...")
	logger.Infof("Subscribing to %d pairs: %s with window size %d", len(productIds), *queryPairs, *vwapWindowSize)

	// Start streaming and handling.
	err = streamHandler.Handle()
	if err != nil {
		logger.Errorf("failed to handle stream data: %v", err)
		return
	}

	// Wait for interrupt signal to gracefully shutdown the process.
	for {
		select {
		case <-interrupt:
			logger.Infoln("Interrupt key signal received, stopping...")
			return
		}
	}
}
