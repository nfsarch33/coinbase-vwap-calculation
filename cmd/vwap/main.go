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
	DefaultWebSocketURL   = "wss://ws-feed.exchange.coinbase.com"
	DefaultLogLevel       = logrus.FatalLevel
	DefaultPairs          = "BTC-USD,ETH-USD,ETH-BTC"
	DefaultVwapWindowSize = 200
)

var (
	queryPairs     = flag.String("pairs", DefaultPairs, "comma separated list of pairs to query")
	verbose        = flag.Bool("verbose", false, "verbose logging")
	wsURL          = flag.String("wsurl", DefaultWebSocketURL, "websocket url")
	vwapWindowSize = flag.Int("window-size", DefaultVwapWindowSize, "vwap window size")
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Set up logger.
	logger := logrus.New()
	logger.SetLevel(DefaultLogLevel)
	if *verbose {
		logger.SetLevel(logrus.TraceLevel)
	}

	productIds := strings.Split(*queryPairs, ",")

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
	streamer := coinbase.NewStreamer(ctx, *wsURL, string(request))
	streamer.SetLogger(logger)

	var streamHandler streaming.StreamDataHandler
	streamHandler = handler.NewStreamDataHandler(*vwapWindowSize, productIds)
	streamHandler.SetLogger(logger)
	streamHandler.SetStreamer(streamer)

	err = streamHandler.Handle()
	if err != nil {
		logger.Errorf("failed to handle stream data: %v", err)
		return
	}

	for {
		select {
		case <-interrupt:
			logger.Infoln("Interrupt key signal received, stopping...")
			return
		}
	}
}
