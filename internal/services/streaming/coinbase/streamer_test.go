package coinbase

import (
	"context"
	"reflect"
	"testing"

	wsclient "bitbucket.org/keynear/coinbase-vwap-calculation/internal/clients/websocket"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming"
	"github.com/sirupsen/logrus"
)

func TestNewStreamer(t *testing.T) {
	type args struct {
		ctx     context.Context
		wsURL   string
		request string
	}

	ctx := context.Background()
	wsUrl := "wss://ws-feed-public.sandbox.exchange.coinbase.com"
	request := "{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[{\"name\":\"matches\",\"product_ids\":[\"BTC-USD\"]}]}"
	logger := logrus.New()
	streamer := NewStreamer(ctx, wsUrl, request)
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
				ctx:     context.Background(),
				wsURL:   wsUrl,
				request: request,
			},
			want: &Streamer{
				ctx:     context.Background(),
				wsURL:   wsUrl,
				client:  wsClient,
				request: request,
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
	client := wsclient.NewClient(context.Background(), "wss://localhost:8080/")
	tests := []struct {
		name   string
		fields fields
		want   *wsclient.Client
	}{
		// Add TestStreamer_GetClient test cases.
		{
			name: "TestStreamer_GetClient",
			fields: fields{
				ctx:               context.Background(),
				wsURL:             "wss://localhost:8080/",
				client:            client,
				request:           "{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[\"ticker\"]}",
				streamDataHandler: nil,
				logger:            logrus.New(),
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

func TestStreamer_GetContext(t *testing.T) {
	type fields struct {
		ctx               context.Context
		wsURL             string
		client            *wsclient.Client
		request           string
		streamDataHandler streaming.StreamDataHandler
		logger            *logrus.Logger
	}
	ctx := context.Background()
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		// Add TestStreamer_GetContext test cases.
		{
			name: "TestStreamer_GetContext",
			fields: fields{
				ctx:               ctx,
				wsURL:             "wss://localhost:8080/",
				client:            nil,
				request:           "{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[\"ticker\"]}",
				streamDataHandler: nil,
				logger:            logrus.New(),
			},
			want: ctx,
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
			if got := s.GetContext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamer_SetLogger(t *testing.T) {
	type fields struct {
		ctx               context.Context
		wsURL             string
		client            *wsclient.Client
		request           string
		streamDataHandler streaming.StreamDataHandler
		logger            *logrus.Logger
	}
	client := wsclient.NewClient(context.Background(), "wss://localhost:8080/")
	type args struct {
		logger *logrus.Logger
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// Add TestStreamer_SetLogger test cases.
		{
			name: "TestStreamer_SetLogger",
			fields: fields{
				ctx:               context.Background(),
				wsURL:             "wss://localhost:8080/",
				client:            client,
				request:           "{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[\"ticker\"]}",
				streamDataHandler: nil,
				logger:            logrus.New(),
			},
			args: args{
				logger: logrus.New(),
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
			s.SetLogger(tt.args.logger)
		})
	}
}

func TestStreamer_SetStreamDataHandler(t *testing.T) {
	type fields struct {
		ctx               context.Context
		wsURL             string
		client            *wsclient.Client
		request           string
		streamDataHandler streaming.StreamDataHandler
		logger            *logrus.Logger
	}
	type args struct {
		streamDataHandler streaming.StreamDataHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// Add TestStreamer_SetStreamDataHandler test cases.
		{
			name: "TestStreamer_SetStreamDataHandler",
			fields: fields{
				ctx:               context.Background(),
				wsURL:             "wss://localhost:8080/",
				client:            nil,
				request:           "{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[\"ticker\"]}",
				streamDataHandler: nil,
				logger:            logrus.New(),
			},
			args: args{
				streamDataHandler: nil,
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
			s.SetStreamDataHandler(tt.args.streamDataHandler)
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
	wsURL := "wss://ws-feed-public.sandbox.exchange.coinbase.com"
	client := wsclient.NewClient(context.Background(), wsURL)
	tests := []struct {
		name   string
		fields fields
	}{
		// Add TestStreamer_Stop test cases.
		{
			name: "TestStreamer_Stop",
			fields: fields{
				ctx:               context.Background(),
				wsURL:             wsURL,
				client:            client,
				request:           "{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[\"ticker\"]}",
				streamDataHandler: nil,
				logger:            logrus.New(),
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
	type fields struct {
		ctx               context.Context
		wsURL             string
		client            *wsclient.Client
		request           string
		streamDataHandler streaming.StreamDataHandler
		logger            *logrus.Logger
	}
	wsURL := "wss://ws-feed.exchange.coinbase.com"
	client := wsclient.NewClient(context.Background(), wsURL)
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
				wsURL:             wsURL,
				client:            client,
				request:           "{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[\"ticker\"]}",
				streamDataHandler: nil,
				logger:            logrus.New(),
			},
			args: args{
				streamFeeds: nil,
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
			defer s.Stop()
			if err := s.Stream(tt.args.streamFeeds); (err != nil) != tt.wantErr {
				t.Errorf("Stream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
