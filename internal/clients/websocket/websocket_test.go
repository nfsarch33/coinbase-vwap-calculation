//go:build all || integration
// +build all integration

package websocket

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	WsURLSandbox = "wss://ws-feed-public.sandbox.exchange.coinbase.com"
	ReqString    = `{"type":"subscribe","product_ids":["BTC-USD"],"channels":{ "name": "matches", "product_ids": ["BTC-USD"]}}`
)

func TestClient_Close(t *testing.T) {
	type fields struct {
		Ctx               context.Context
		Conn              *websocket.Conn
		WebsocketDialer   *websocket.Dialer
		URL               string
		ConnectionOptions ConnOptions
		RequestHeader     http.Header
		OnConnected       func(client Client)
		OnReceivingMsg    func(message string, client Client)
		OnConnectError    func(err error, client Client)
		OnDisconnected    func(err error, client Client)
		IsConnected       bool
		Timeout           time.Duration
		sendMu            *sync.Mutex
		receiveMu         *sync.Mutex
		logger            *logrus.Logger
	}
	logger := logrus.New()
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "TestClient_Close",
			fields: fields{
				Ctx:             context.Background(),
				Conn:            nil,
				WebsocketDialer: &websocket.Dialer{},
				URL:             WsURLSandbox,
				ConnectionOptions: ConnOptions{
					UseCompression: false,
					UseSSL:         true,
				},
				RequestHeader:  http.Header{},
				OnConnected:    nil,
				OnReceivingMsg: nil,
				OnConnectError: nil,
				OnDisconnected: nil,
				IsConnected:    true,
				Timeout:        0,
				sendMu:         &sync.Mutex{},
				receiveMu:      &sync.Mutex{},
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Ctx:               tt.fields.Ctx,
				Conn:              tt.fields.Conn,
				WebsocketDialer:   tt.fields.WebsocketDialer,
				URL:               tt.fields.URL,
				ConnectionOptions: tt.fields.ConnectionOptions,
				RequestHeader:     tt.fields.RequestHeader,
				OnConnected:       tt.fields.OnConnected,
				OnReceivingMsg:    tt.fields.OnReceivingMsg,
				OnConnectError:    tt.fields.OnConnectError,
				OnDisconnected:    tt.fields.OnDisconnected,
				IsConnected:       tt.fields.IsConnected,
				Timeout:           tt.fields.Timeout,
				sendMu:            tt.fields.sendMu,
				receiveMu:         tt.fields.receiveMu,
				logger:            tt.fields.logger,
			}
			c.OnDisconnected = func(err error, socket Client) {
				if err != nil {
					logger.Errorf("Received disconnect error %s", err)
				} else {
					logger.Infoln("Disconnected from server")
				}
			}

			c.Close()

			if c.Conn != nil {
				t.Errorf("Client.Close() Conn = %v, want %v", c.Conn, nil)
			}
		})
	}
}

func TestClient_Connect(t *testing.T) {
	type fields struct {
		Ctx               context.Context
		Conn              *websocket.Conn
		WebsocketDialer   *websocket.Dialer
		URL               string
		ConnectionOptions ConnOptions
		RequestHeader     http.Header
		OnConnected       func(client Client)
		OnReceivingMsg    func(message string, client Client)
		OnConnectError    func(err error, client Client)
		OnDisconnected    func(err error, client Client)
		IsConnected       bool
		Timeout           time.Duration
		sendMu            *sync.Mutex
		receiveMu         *sync.Mutex
		logger            *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// Test wsclinet connect function.
		{
			name: "TestClient_Connect",
			fields: fields{
				Ctx:             context.Background(),
				Conn:            nil,
				WebsocketDialer: &websocket.Dialer{},
				URL:             WsURLSandbox,
				ConnectionOptions: ConnOptions{
					UseCompression: false,
					UseSSL:         true,
				},
				RequestHeader:  http.Header{},
				OnConnected:    nil,
				OnReceivingMsg: nil,
				OnConnectError: nil,
				OnDisconnected: nil,
				IsConnected:    true,
				Timeout:        0,
				sendMu:         &sync.Mutex{},
				receiveMu:      &sync.Mutex{},
				logger:         logrus.New(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Ctx:               tt.fields.Ctx,
				Conn:              tt.fields.Conn,
				WebsocketDialer:   tt.fields.WebsocketDialer,
				URL:               tt.fields.URL,
				ConnectionOptions: tt.fields.ConnectionOptions,
				RequestHeader:     tt.fields.RequestHeader,
				OnConnected:       tt.fields.OnConnected,
				OnReceivingMsg:    tt.fields.OnReceivingMsg,
				OnConnectError:    tt.fields.OnConnectError,
				OnDisconnected:    tt.fields.OnDisconnected,
				IsConnected:       tt.fields.IsConnected,
				Timeout:           tt.fields.Timeout,
				sendMu:            tt.fields.sendMu,
				receiveMu:         tt.fields.receiveMu,
				logger:            tt.fields.logger,
			}
			if err := c.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !c.IsConnected {
				t.Errorf("Connect() IsConnected = %v, want %v", c.IsConnected, true)
			}
		})
	}
}

func TestClient_SendRequest(t *testing.T) {
	type fields struct {
		Ctx               context.Context
		Conn              *websocket.Conn
		WebsocketDialer   *websocket.Dialer
		URL               string
		ConnectionOptions ConnOptions
		RequestHeader     http.Header
		OnConnected       func(client Client)
		OnReceivingMsg    func(message string, client Client)
		OnConnectError    func(err error, client Client)
		OnDisconnected    func(err error, client Client)
		IsConnected       bool
		Timeout           time.Duration
		sendMu            *sync.Mutex
		receiveMu         *sync.Mutex
		logger            *logrus.Logger
	}
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Test client send request function.
		{
			name: "TestClient_SendRequest",
			fields: fields{
				Ctx:             context.Background(),
				Conn:            nil,
				WebsocketDialer: &websocket.Dialer{},
				URL:             WsURLSandbox,
				ConnectionOptions: ConnOptions{
					UseCompression: false,
					UseSSL:         true,
				},
				RequestHeader:  http.Header{},
				OnConnected:    nil,
				OnReceivingMsg: nil,
				OnConnectError: nil,
				OnDisconnected: nil,
				IsConnected:    false,
				Timeout:        0,
				sendMu:         &sync.Mutex{},
				receiveMu:      &sync.Mutex{},
				logger:         logrus.New(),
			},
			args: args{
				message: ReqString,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Ctx:               tt.fields.Ctx,
				Conn:              tt.fields.Conn,
				WebsocketDialer:   tt.fields.WebsocketDialer,
				URL:               tt.fields.URL,
				ConnectionOptions: tt.fields.ConnectionOptions,
				RequestHeader:     tt.fields.RequestHeader,
				OnConnected:       tt.fields.OnConnected,
				OnReceivingMsg:    tt.fields.OnReceivingMsg,
				OnConnectError:    tt.fields.OnConnectError,
				OnDisconnected:    tt.fields.OnDisconnected,
				IsConnected:       tt.fields.IsConnected,
				Timeout:           tt.fields.Timeout,
				sendMu:            tt.fields.sendMu,
				receiveMu:         tt.fields.receiveMu,
				logger:            tt.fields.logger,
			}
			err := c.Connect()
			if err != nil {
				t.Errorf("Connect() error before SendRequest() error = %v", err)
			}
			if err := c.SendRequest(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("SendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			c.OnReceivingMsg = func(message string, socket Client) {
				t.Logf("Received message: %s", message)
				if message == "" {
					t.Errorf("SendRequest() OnReceivingMsg = %v, want %v", message, ReqString)
				}
			}
		})
	}
}

func TestClient_send(t *testing.T) {
	type fields struct {
		Ctx               context.Context
		Conn              *websocket.Conn
		WebsocketDialer   *websocket.Dialer
		URL               string
		ConnectionOptions ConnOptions
		RequestHeader     http.Header
		OnConnected       func(client Client)
		OnReceivingMsg    func(message string, client Client)
		OnConnectError    func(err error, client Client)
		OnDisconnected    func(err error, client Client)
		IsConnected       bool
		Timeout           time.Duration
		sendMu            *sync.Mutex
		receiveMu         *sync.Mutex
		logger            *logrus.Logger
	}
	type args struct {
		messageType int
		data        []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Add TestClient_send test cases.
		{
			name: "TestClient_send",
			fields: fields{
				Ctx:             context.Background(),
				Conn:            nil,
				WebsocketDialer: &websocket.Dialer{},
				URL:             WsURLSandbox,
				ConnectionOptions: ConnOptions{
					UseCompression: false,
					UseSSL:         true,
				},
				RequestHeader:  http.Header{},
				OnConnected:    nil,
				OnReceivingMsg: nil,
				OnConnectError: nil,
				OnDisconnected: nil,
				IsConnected:    false,
				Timeout:        0,
				sendMu:         &sync.Mutex{},
				receiveMu:      &sync.Mutex{},
				logger:         logrus.New(),
			},
			args: args{
				messageType: websocket.TextMessage,
				data:        []byte(ReqString),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Ctx:               tt.fields.Ctx,
				Conn:              tt.fields.Conn,
				WebsocketDialer:   tt.fields.WebsocketDialer,
				URL:               tt.fields.URL,
				ConnectionOptions: tt.fields.ConnectionOptions,
				RequestHeader:     tt.fields.RequestHeader,
				OnConnected:       tt.fields.OnConnected,
				OnReceivingMsg:    tt.fields.OnReceivingMsg,
				OnConnectError:    tt.fields.OnConnectError,
				OnDisconnected:    tt.fields.OnDisconnected,
				IsConnected:       tt.fields.IsConnected,
				Timeout:           tt.fields.Timeout,
				sendMu:            tt.fields.sendMu,
				receiveMu:         tt.fields.receiveMu,
				logger:            tt.fields.logger,
			}
			err := c.Connect()
			if err != nil {
				t.Errorf("Connect() error before send() error = %v", err)
			}
			if err := c.send(tt.args.messageType, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		ctx   context.Context
		wsURL string
	}
	logger := logrus.New()
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// Add TestNewClient function test cases.
		{
			name: "TestNewClient",
			args: args{
				ctx:   context.Background(),
				wsURL: WsURLSandbox,
			},
			want: &Client{
				Ctx:           context.Background(),
				URL:           WsURLSandbox,
				RequestHeader: http.Header{},
				ConnectionOptions: ConnOptions{
					UseCompression: false,
					UseSSL:         true,
				},
				WebsocketDialer: &websocket.Dialer{},
				Timeout:         0,
				sendMu:          &sync.Mutex{},
				receiveMu:       &sync.Mutex{},
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.args.ctx, tt.args.wsURL)
			got.SetLogger(logger)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
