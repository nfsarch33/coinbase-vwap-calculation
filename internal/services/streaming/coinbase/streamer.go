package coinbase

import (
	wsclient "bitbucket.org/keynear/coinbase-vwap-calculation/internal/clients/websocket"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

const (
	FeedTypeMatch          = "match"
	FeedTypeSubscribeError = "error"
	FeedTypeLastMatch      = "last_match"
	FeedTypeLevel2Snapshot = "l2update"
	FeedTypeTicker         = "ticker"
)

// Streamer is a streaming service for Coinbase. It implements the streaming.Streamer interface.
// It consists of a websocket client and a message handler streamDataHandler.
type Streamer struct {
	ctx               context.Context
	wsURL             string
	client            *wsclient.Client
	request           string
	streamDataHandler streaming.StreamDataHandler
	logger            *logrus.Logger
}

func NewStreamer(ctx context.Context, wsURL string, request string) *Streamer {
	return &Streamer{
		ctx:     ctx,
		wsURL:   wsURL,
		client:  wsclient.NewClient(ctx, wsURL),
		request: request,
		logger:  logrus.New(),
	}
}

func (s *Streamer) SetLogger(logger *logrus.Logger) {
	s.logger = logger
	s.client.SetLogger(logger)
}

func (s *Streamer) SetStreamDataHandler(streamDataHandler streaming.StreamDataHandler) {
	s.streamDataHandler = streamDataHandler
}

func (s *Streamer) GetClient() *wsclient.Client {
	return s.client
}

func (s *Streamer) GetContext() context.Context {
	return s.ctx
}

// Stream starts the process of subscribing to a channel and streaming feeds from Coinbase, the datapoint is passed to
// a streamFeeds channel that can be further passed to the stream data handler.
func (s *Streamer) Stream(
	streamFeeds chan interface{},
) error {
	client := s.client

	var cancel context.CancelFunc
	s.ctx, cancel = context.WithCancel(s.GetContext())

	client.OnConnected = func(socket wsclient.Client) {
		s.logger.Infoln("Connected to coinbase server.")
	}

	client.OnConnectError = func(err error, socket wsclient.Client) {
		s.logger.Infoln("Received connect error ", err)
	}

	client.OnReceivingMsg = func(message string, socket wsclient.Client) {
		var m = Feed{}

		err := json.Unmarshal([]byte(message), &m)
		if err != nil {
			s.logger.Errorf("Error unmarshalling message %s", err)
			return
		}

		// Stop on subscribe errors.
		if m.Type == FeedTypeSubscribeError {
			s.logger.Errorf("Received subscribe error type: %v error: %v", m.Type, m)
			s.logger.Errorf("Reason: %v", m.Reason)
			cancel()
			return
		}

		go func() {
			// This is to prevent race condition upon connection error or closed connection.
			if !client.IsConnected || client.OnConnected == nil {
				cancel()
				return
			}

			if m.Type == FeedTypeMatch || m.Type == FeedTypeLastMatch {
				streamFeeds <- m
			}
		}()
	}

	client.OnDisconnected = func(err error, socket wsclient.Client) {
		if err != nil {
			s.logger.Errorf("Received disconnect error %s", err)
		} else {
			s.logger.Infoln("Disconnected from server")
		}
	}

	if !client.IsConnected {
		err := client.Connect()
		if err != nil {
			s.logger.Errorf("Error connecting to server %s", err)

			return err
		}
	}

	if client.IsConnected {
		err := client.SendRequest(s.request)
		if err != nil {
			s.logger.Errorf("Error sending request %s", err)

			return err
		}
	}

	return nil
}

func (s *Streamer) Stop() {
	var cancel context.CancelFunc
	s.ctx, cancel = context.WithCancel(s.ctx)
	defer cancel()

	if s.client.IsConnected {
		s.client.Close()
	}
}
