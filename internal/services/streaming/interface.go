package streaming

import (
	"context"

	wsclient "bitbucket.org/keynear/coinbase-vwap-calculation/internal/clients/websocket"
	"github.com/sirupsen/logrus"
)

// StreamDataHandler is the interface for implementing incoming data processing handler.
type StreamDataHandler interface {
	SetStreamer(streamer Streamer)
	SetLogger(logger *logrus.Logger)
	Handle() error
}

// Streamer is the interface for implementing streaming feeds.
type Streamer interface {
	GetContext() context.Context
	GetClient() *wsclient.Client
	SetLogger(logger *logrus.Logger)
	Stream(streamFeeds chan interface{}) error
}
