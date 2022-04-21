//go:build all || integration
// +build all integration

package coinbase

import (
	wsclient "bitbucket.org/keynear/coinbase-vwap-calculation/internal/clients/websocket"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.uber.org/goleak"
	"reflect"
	"testing"
)

const (
	WsURLSandbox = "wss://ws-feed.exchange.coinbase.com"
	ReqString    = `
					{
						"type": "subscribe",
						"product_ids": [
							"BTC-USD"
						],
						"channels": [{
							"name": "matches",
							"product_ids": [
								"BTC-USD"
							]
						}]
					}
					`
)

func TestNewStreamer(t *testing.T) {
	type args struct {
		ctx     context.Context
		wsURL   string
		request string
	}

	ctx := context.Background()
	logger := logrus.New()

	streamer := NewStreamer(ctx, WsURLSandbox, ReqString)
	wsClient := streamer.GetClient()
	tests := []struct {
		name string
		args args
		want *Streamer
	}{
		// Add TestNewStreamer test cases.
		{
			name: "TestNewStreamer",
			args: args{
				ctx:     ctx,
				wsURL:   WsURLSandbox,
				request: ReqString,
			},
			want: &Streamer{
				ctx:     ctx,
				wsURL:   WsURLSandbox,
				client:  wsClient,
				request: ReqString,
				logger:  logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := streamer
			got.SetLogger(logger)
			if reflect.DeepEqual(got, tt.want) != true {
				t.Errorf("NewStreamer() = %v, \n want %v", got, tt.want)
			}
		})
	}
}

func TestStreamer_GetClient(t *testing.T) {
	type fields struct {
		ctx               context.Context
		wsURL             string
		client            *wsclient.Client
		request           string
		streamDataHandler streaming.StreamDataHandler
		logger            *logrus.Logger
	}

	ctx := context.Background()
	logger := logrus.New()

	client := wsclient.NewClient(ctx, WsURLSandbox)
	tests := []struct {
		name   string
		fields fields
		want   *wsclient.Client
	}{
		// Add TestStreamer_GetClient test cases.
		{
			name: "TestStreamer_GetClient",
			fields: fields{
				ctx:               ctx,
				wsURL:             WsURLSandbox,
				client:            client,
				request:           ReqString,
				streamDataHandler: nil,
				logger:            logger,
			},
			want: client,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Streamer{
				ctx:               tt.fields.ctx,
				wsURL:             tt.fields.wsURL,
				client:            tt.fields.client,
				request:           tt.fields.request,
				streamDataHandler: tt.fields.streamDataHandler,
				logger:            tt.fields.logger,
			}
			if got := s.GetClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func TestStreamer_Stop(t *testing.T) {
	type fields struct {
		ctx               context.Context
		wsURL             string
		client            *wsclient.Client
		request           string
		streamDataHandler streaming.StreamDataHandler
		logger            *logrus.Logger
	}

	ctx := context.Background()
	logger := logrus.New()

	client := wsclient.NewClient(ctx, WsURLSandbox)
	tests := []struct {
		name   string
		fields fields
	}{
		// Add TestStreamer_Stop test cases.
		{
			name: "TestStreamer_Stop",
			fields: fields{
				ctx:               ctx,
				wsURL:             WsURLSandbox,
				client:            client,
				request:           ReqString,
				streamDataHandler: nil,
				logger:            logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Streamer{
				ctx:               tt.fields.ctx,
				wsURL:             tt.fields.wsURL,
				client:            tt.fields.client,
				request:           tt.fields.request,
				streamDataHandler: tt.fields.streamDataHandler,
				logger:            tt.fields.logger,
			}
			s.Stop()
			if s.GetClient().IsConnected {
				t.Errorf(
					"Streamer.Stop() client.IsConnected = %v, want %v",
					s.GetClient().IsConnected,
					false,
				)
			}
		})
	}
}

func TestStreamer_Stream(t *testing.T) {
	defer goleak.VerifyNone(t)
	type fields struct {
		ctx               context.Context
		wsURL             string
		client            *wsclient.Client
		request           string
		streamDataHandler streaming.StreamDataHandler
		logger            *logrus.Logger
	}

	ctx := context.Background()
	logger := logrus.New()

	client := wsclient.NewClient(ctx, WsURLSandbox)
	type args struct {
		streamFeeds chan interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Add TestStreamer_Stream test cases.
		{
			name: "TestStreamer_Stream",
			fields: fields{
				ctx:               context.Background(),
				wsURL:             WsURLSandbox,
				client:            client,
				request:           ReqString,
				streamDataHandler: nil,
				logger:            logger,
			},
			args: args{
				streamFeeds: make(chan interface{}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Streamer{
				ctx:               tt.fields.ctx,
				wsURL:             tt.fields.wsURL,
				client:            tt.fields.client,
				request:           tt.fields.request,
				streamDataHandler: tt.fields.streamDataHandler,
				logger:            tt.fields.logger,
			}
			defer close(tt.args.streamFeeds)
			defer s.Stop()
			if err := s.Stream(tt.args.streamFeeds); (err != nil) != tt.wantErr {
				t.Errorf("Stream() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Test receiving data from stream, test max 10 samples.
			count := 0
			maxDataPoints := 3

			for {
				if count >= maxDataPoints {
					break
				}

				select {
				case <-s.ctx.Done():
					t.Errorf("Streamer.Stream() context.Done()")
					return
				case feed := <-tt.args.streamFeeds:

					f, err := interfaceToFeedStruct(feed)
					if err != nil {
						t.Errorf("Streamer.Stream() InterfaceToFeedStruct error = %v, wantErr %v", err, tt.wantErr)
					}
					t.Logf("Streamer.Stream() feed = %v", f)

					if f.Type == "error" {
						t.Errorf("Streamer.Stream() feed.Type = %v, want %v", f.Type, "error")
					}

					switch f.Type {
					case "match":
						t.Logf("Streamer.Stream() match = %v", f)
						if f.ProductID != "BTC-USD" {
							t.Errorf("Streamer.Stream() feed.ProductID = %v, want %v", f.ProductID, "BTC-USD")
						}
					case "last_match":
						t.Logf("Streamer.Stream() last_match = %v", f)
						if f.ProductID != "BTC-USD" {
							t.Errorf("Streamer.Stream() feed.ProductID = %v, want %v", f.ProductID, "BTC-USD")
						}
					default:
						t.Errorf("Streamer.Stream() feed.Type = %v, want %v", f.Type, "match")
					}

					count++
				}
			}
		})
	}
}

var interfaceToFeedStruct = func(anyData interface{}) (Feed, error) {
	bytes, err := json.Marshal(anyData)
	if err != nil {
		return Feed{}, err
	}

	feed := Feed{}

	err = json.Unmarshal(bytes, &feed)
	if err != nil {
		return Feed{}, err
	}

	return feed, nil
}
