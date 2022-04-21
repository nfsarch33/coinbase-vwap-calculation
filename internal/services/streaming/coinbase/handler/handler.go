package handler

import (
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming/coinbase"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/vwap"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

// CoinbaseSteamDataHandler is the implementation of the streaming.DataHandler interface.
// It is used to handle the incoming data from the Coinbase streaming API wrapped by streamer.
type CoinbaseSteamDataHandler struct {
	vwapMaxSize         int
	vwapPairs           []string
	vwapData            map[string]*vwap.SlidingWindow
	MessagePipelineFunc func(s *vwap.SlidingWindow) error
	streamer            streaming.Streamer
	logger              *logrus.Logger
}

func NewStreamDataHandler(maxSize int, pairs []string) *CoinbaseSteamDataHandler {
	return &CoinbaseSteamDataHandler{
		vwapMaxSize: maxSize,
		vwapPairs:   pairs,
		vwapData:    make(map[string]*vwap.SlidingWindow),
		logger:      logrus.New(),
	}
}

func (h *CoinbaseSteamDataHandler) SetLogger(logger *logrus.Logger) {
	h.logger = logger
}

func (h *CoinbaseSteamDataHandler) SetStreamer(streamer streaming.Streamer) {
	h.streamer = streamer
}

func (h *CoinbaseSteamDataHandler) GetStreamer() streaming.Streamer {
	return h.streamer
}

// SetMessageBlockerFunc SetMessagePipelineFunc sets the function that will be called when a new message is received.
func (h *CoinbaseSteamDataHandler) SetMessageBlockerFunc(
	msgBlockerFunc func(c *vwap.SlidingWindow) error,
) {
	h.MessagePipelineFunc = msgBlockerFunc
}

// Handle handles the incoming data from the streamer and pipes it to a MessagePipelineFunc
// that can be implemented later.
func (h *CoinbaseSteamDataHandler) Handle() error {
	s := h.streamer
	streamFeeds := make(chan interface{})
	ctx := s.GetContext()

	err := s.Stream(streamFeeds)
	if err != nil {
		h.logger.Errorf("Error starting stream %s", err)
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(streamFeeds)
				s.GetClient().Close()
				return
			case feed := <-streamFeeds:
				f, err := InterfaceToFeedStruct(feed)

				if err != nil {
					h.logger.Errorf("Error converting interface to feed struct %s", err)
					continue
				}

				dataPoint := vwap.DataPoint{
					Type:      f.Type,
					Size:      f.Size,
					Price:     f.Price,
					ProductID: f.ProductID,
				}

				err = h.processVwapData(dataPoint)
				if err != nil {
					h.logger.Errorf("Error processing vwap data %s", err)
					continue
				}

				// TODO: Implement message pipeline function to send it to the message blocker or DB.
				if h.MessagePipelineFunc != nil {
					err := h.MessagePipelineFunc(h.vwapData[dataPoint.ProductID])
					if err != nil {
						h.logger.Errorf("Error processing vwap data %s", err)
						continue
					}
				}
			}
		}
	}()

	return nil
}

// processVwapData processes the incoming feed data and updates the vwap data property.
func (h *CoinbaseSteamDataHandler) processVwapData(dataPoint vwap.DataPoint) error {
	if _, ok := h.vwapData[dataPoint.ProductID]; !ok {
		h.vwapData[dataPoint.ProductID] = vwap.NewSlidingWindow(h.vwapMaxSize, dataPoint.ProductID)
	}

	h.vwapData[dataPoint.ProductID].Add(dataPoint)

	fmt.Printf(
		"Windows Size: %v\t%v:%v\n",
		h.vwapData[dataPoint.ProductID].Size(),
		dataPoint.ProductID,
		h.vwapData[dataPoint.ProductID].GetCalculator().Avg().String())

	return nil
}

func InterfaceToFeedStruct(anyData interface{}) (coinbase.Feed, error) {
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
