//go:build all || integration
// +build all integration

package handler

import (
	"context"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming/coinbase"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/vwap"
	"github.com/sirupsen/logrus"
)

const (
	WsURLSandbox = "wss://ws-feed-public.sandbox.exchange.coinbase.com"
	ReqString    = `
					{
						"type": "subscribe",
						"product_ids": [
							"BTC-USD",
							"ETH-USD",
							"ETH-BTC"
						],
						"channels": [{
							"name": "matches",
							"product_ids": [
								"BTC-USD",
								"ETH-USD",
								"ETH-BTC"
							]
						}]
					}
					`
)

var (
	ctx       = context.Background()
	testPair  = []string{"BTC-USD"}
	testPairs = []string{"BTC-USD", "ETH-USD", "LTC-USD"}
	logger    = logrus.New()
)

func TestCoinbaseSteamDataHandler_Handle(t *testing.T) {
	type fields struct {
		vwapMaxSize         int
		vwapPairs           []string
		vwapData            map[string]*vwap.SlidingWindow
		messagePipelineFunc func(s *vwap.SlidingWindow) error
		streamer            streaming.Streamer
		logger              *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// Add TestCoinbaseSteamDataHandler_Handle test cases.
		{
			name: "TestCoinbaseSteamDataHandler_Handle",
			fields: fields{
				vwapMaxSize:         10,
				vwapPairs:           []string{"BTC-USD", "ETH-USD", "ETH-BTC"},
				vwapData:            make(map[string]*vwap.SlidingWindow),
				messagePipelineFunc: func(s *vwap.SlidingWindow) error { return nil },
				streamer: coinbase.NewStreamer(
					ctx,
					WsURLSandbox,
					ReqString,
				),
				logger: logger,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CoinbaseSteamDataHandler{
				vwapMaxSize:         tt.fields.vwapMaxSize,
				vwapPairs:           tt.fields.vwapPairs,
				vwapData:            tt.fields.vwapData,
				messagePipelineFunc: tt.fields.messagePipelineFunc,
				streamer:            tt.fields.streamer,
				logger:              tt.fields.logger,
			}
			if err := h.Handle(); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoinbaseSteamDataHandler_processVwapData(t *testing.T) {
	type fields struct {
		vwapMaxSize         int
		vwapPairs           []string
		vwapData            map[string]*vwap.SlidingWindow
		messagePipelineFunc func(s *vwap.SlidingWindow) error
		streamer            streaming.Streamer
		logger              *logrus.Logger
	}
	type args struct {
		dataPoint vwap.DataPoint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Add TestCoinbaseSteamDataHandler_processVwapData test cases.
		{
			name: "TestCoinbaseSteamDataHandler_processVwapData",
			fields: fields{
				vwapMaxSize:         10,
				vwapPairs:           []string{"BTC-USD"},
				vwapData:            make(map[string]*vwap.SlidingWindow),
				messagePipelineFunc: func(s *vwap.SlidingWindow) error { return nil },
				streamer: coinbase.NewStreamer(
					ctx,
					WsURLSandbox,
					ReqString,
				),
				logger: logger,
			},
			args: args{
				dataPoint: vwap.DataPoint{
					Type:      "match",
					Price:     big.NewFloat(1.0),
					Size:      big.NewFloat(1.0),
					ProductID: "BTC-USD",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CoinbaseSteamDataHandler{
				vwapMaxSize:         tt.fields.vwapMaxSize,
				vwapPairs:           tt.fields.vwapPairs,
				vwapData:            tt.fields.vwapData,
				messagePipelineFunc: tt.fields.messagePipelineFunc,
				streamer:            tt.fields.streamer,
				logger:              tt.fields.logger,
			}
			if err := h.processVwapData(tt.args.dataPoint); (err != nil) != tt.wantErr {
				t.Errorf("processVwapData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewStreamDataHandler(t *testing.T) {
	logger := logger
	type args struct {
		maxSize int
		pairs   []string
	}
	tests := []struct {
		name string
		args args
		want *CoinbaseSteamDataHandler
	}{
		// Add TestNewStreamDataHandler test cases.
		{
			name: "TestNewStreamDataHandler",
			args: args{
				maxSize: 5,
				pairs:   testPairs,
			},
			want: &CoinbaseSteamDataHandler{
				vwapMaxSize: 5,
				vwapPairs:   testPairs,
				vwapData:    make(map[string]*vwap.SlidingWindow),
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStreamDataHandler(tt.args.maxSize, tt.args.pairs)
			got.SetLogger(logger)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStreamDataHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interfaceToFeedStruct(t *testing.T) {
	type args struct {
		anyData interface{}
	}

	str := `{
				"type": "match",
				"trade_id": 255888086,
				"maker_order_id": "47db4060-895d-4b7f-b46e-ba9742a876f3",
				"taker_order_id": "e9ad0749-9a11-494f-81c1-99af598cdee6",
				"side": "sell",
				"size": "0.01",
				"price": "3005.71",
				"product_id": "ETH-USD",
				"sequence": 28004452801,
				"time": "2022-04-13T13:08:16.502461Z"
			}`

	var data interface{}
	feed := coinbase.Feed{}

	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		t.Errorf("Unmarshal failed %v", err)
	}

	err = json.Unmarshal([]byte(str), &feed)
	if err != nil {
		t.Errorf("Unmarshal failed %v", err)
	}

	tests := []struct {
		name    string
		args    args
		want    coinbase.Feed
		wantErr bool
	}{
		// Add Test_interfaceToFeedStruct test cases.
		{
			name: "Test_interfaceToFeedStruct",
			args: args{
				anyData: data,
			},
			want:    feed,
			wantErr: false,
		},
	}

	InterfaceToFeedStruct := func(anyData interface{}) (coinbase.Feed, error) {
		bytes, err := json.Marshal(anyData)
		if err != nil {
			return coinbase.Feed{}, err
		}

		feed := coinbase.Feed{}

		err = json.Unmarshal(bytes, &feed)
		if err != nil {
			return coinbase.Feed{}, err
		}

		return feed, nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InterfaceToFeedStruct(tt.args.anyData)
			if (err != nil) != tt.wantErr {
				t.Errorf("interfaceToFeedStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("interfaceToFeedStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
