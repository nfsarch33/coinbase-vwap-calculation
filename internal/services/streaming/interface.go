package streaming

import (
	"context"

	wsclient "bitbucket.org/keynear/coinbase-vwap-calculation/internal/clients/websocket"
	"github.com/sirupsen/logrus"
)

type StreamDataHandler interface {
	SetStreamer(streamer Streamer)
	SetLogger(logger *logrus.Logger)
	Handle() error
}

type Streamer interface {
	GetContext() context.Context
	GetClient() *wsclient.Client
	SetLogger(logger *logrus.Logger)
	Stream(streamFeeds chan interface{}) error
}
