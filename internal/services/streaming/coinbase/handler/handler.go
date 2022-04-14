package handler

import (
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/services/streaming/coinbase"
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/vwap"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type CoinbaseSteamDataHandler struct {
	vwapMaxSize         int
	vwapPairs           []string
	vwapData            map[string]*vwap.SlidingWindow
	messagePipelineFunc func(s *vwap.SlidingWindow) error
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

func (h *CoinbaseSteamDataHandler) SetMessageBlockerFunc(
	msgBlockerFunc func(c *vwap.SlidingWindow) error,
) {
	h.messagePipelineFunc = msgBlockerFunc
}

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
				//h.logger.Warning("Streaming context timed out or canceled.")
				s.GetClient().Close()
				return
			case feed := <-streamFeeds:
				f, err := interfaceToFeedStruct(feed)
				//fmt.Printf(f.Type, f.Price, f.Size, f.Time)
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
				if h.messagePipelineFunc != nil {
					err := h.messagePipelineFunc(h.vwapData[dataPoint.ProductID])
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

func (h *CoinbaseSteamDataHandler) processVwapData(dataPoint vwap.DataPoint) error {
	if _, ok := h.vwapData[dataPoint.ProductID]; !ok {
		h.vwapData[dataPoint.ProductID] = vwap.NewSlidingWindow(h.vwapMaxSize, dataPoint.ProductID)
	}

	h.vwapData[dataPoint.ProductID].Add(dataPoint)

	fmt.Printf(
		"%v\t:%v\n",
		dataPoint.ProductID,
		h.vwapData[dataPoint.ProductID].GetCalculator().Avg().String())

	return nil
}

func interfaceToFeedStruct(anyData interface{}) (coinbase.Feed, error) {
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
