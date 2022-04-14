package websocket

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client is a generic websocket client that build on top of gorilla/websocket package. It implements the event driven
// pattern and provides a set of simple functions to receive messages.
type Client struct {
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

type ConnOptions struct {
	UseCompression bool
	UseSSL         bool
	Proxy          func(*http.Request) (*url.URL, error)
	SubProtocols   []string
}

func NewClient(ctx context.Context, wsURL string) *Client {
	return &Client{
		Ctx:           ctx,
		URL:           wsURL,
		RequestHeader: http.Header{},
		ConnectionOptions: ConnOptions{
			UseCompression: false,
			UseSSL:         true,
		},
		WebsocketDialer: &websocket.Dialer{},
		Timeout:         0,
		sendMu:          &sync.Mutex{},
		receiveMu:       &sync.Mutex{},
		logger:          logrus.New(),
	}
}

func (c *Client) SetLogger(logger *logrus.Logger) {
	c.logger = logger
}

func (c *Client) setConnectionOptions() {
	c.WebsocketDialer.EnableCompression = c.ConnectionOptions.UseCompression
	c.WebsocketDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: c.ConnectionOptions.UseSSL}
	c.WebsocketDialer.Proxy = c.ConnectionOptions.Proxy
	c.WebsocketDialer.Subprotocols = c.ConnectionOptions.SubProtocols
}

// Connect connects to the websocket server with a given request message, it then pipes the incoming messages to the
// OnReceivingMsg callback receiver for further processing.
func (c *Client) Connect() error {
	var (
		err    error
		resp   *http.Response
		logger = c.logger
	)

	c.setConnectionOptions()

	// Connect to the websocket server.
	c.Conn, resp, err = c.WebsocketDialer.Dial(c.URL, c.RequestHeader)

	if err != nil {
		logger.Errorf("Error connecting to websocket: %s", err)
		if resp != nil {
			logger.Errorf("HTTP Response %d status: %s", resp.StatusCode, resp.Status)
		}

		c.IsConnected = false

		// Set the OnConnectError callback.
		if c.OnConnectError != nil {
			c.OnConnectError(err, *c)
		}

		return err
	}

	if c.OnConnected != nil {
		c.IsConnected = true
		c.OnConnected(*c)
	}

	logger.Infoln("Connected to server")

	// Set the OnDisconnected callback.
	defaultCloseHandler := c.Conn.CloseHandler()
	c.Conn.SetCloseHandler(func(code int, text string) error {
		result := defaultCloseHandler(code, text)
		logger.Warning("Disconnected from server ", result)
		if c.OnDisconnected != nil {
			c.IsConnected = false
			c.OnDisconnected(errors.New(text), *c)
		}

		return result
	})

	go func() {
		for {
			c.receiveMu.Lock()
			if c.Timeout != 0 {
				err := c.Conn.SetReadDeadline(time.Now().Add(c.Timeout))
				if err != nil {
					logger.Errorf("Error setting read deadline: %s", err)

					return
				}
			}

			messageType, message, err := c.Conn.ReadMessage()
			if messageType == -1 {
				logger.Errorf("Error reading message: %s", err)
				c.Close()

				return
			}

			if err != nil {
				logger.Errorf("read: %s", err)
				if c.OnDisconnected != nil {
					c.IsConnected = false
				}

				c.OnDisconnected(err, *c)

				return
			}

			// Pipe the response message to the OnReceivingMsg callback receiver.
			if websocket.TextMessage == messageType && c.OnReceivingMsg != nil {
				c.OnReceivingMsg(string(message), *c)
			}

			c.receiveMu.Unlock()
		}
	}()

	return nil
}

// Send is a proxy function of WriteMessage, it sends a message to the websocket server.
func (c *Client) send(messageType int, data []byte) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	return c.Conn.WriteMessage(messageType, data)
}

func (c *Client) SendRequest(message string) error {
	err := c.send(websocket.TextMessage, []byte(message))
	if err != nil {
		c.logger.Errorf("write: %s", err)

		return err
	}

	return nil
}

// Close closes the websocket connection.
func (c *Client) Close() {
	logger := c.logger

	if !c.IsConnected || c.Conn == nil {
		return
	}

	err := c.send(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		logger.Errorf("write close: %s", err)

		return
	}

	err = c.Conn.Close()
	if err != nil {
		logger.Errorf("close: %s", err)

		return
	}

	// Set the OnDisconnected callback.
	if c.OnDisconnected != nil {
		c.IsConnected = false
		c.OnDisconnected(err, *c)
	}
}
