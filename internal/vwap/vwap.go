package vwap

import (
	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/vwap/utils"
	"math/big"
	"sync"
)

var mux sync.Mutex

// SlidingWindow is a struct that holds the sliding window that contains a set of datapoints.
// A VolumeWeightedAveragePriceCalculator is attached to hold the total volume and the total price.
type SlidingWindow struct {
	currencyPair string
	dataPoints   []DataPoint
	windowSize   int
	calculator   utils.VolumeWeightedAveragePriceCalculator
}

type DataPoint struct {
	Type      string
	Size      *big.Float
	Price     *big.Float
	ProductID string
}

func NewSlidingWindow(maxSize int, currencyPair string) *SlidingWindow {
	return &SlidingWindow{
		currencyPair: currencyPair,
		dataPoints:   make([]DataPoint, 0),
		windowSize:   maxSize,
		calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
	}
}

func (sw *SlidingWindow) SetSize(maxSize int) {
	sw.windowSize = maxSize
}

func (sw *SlidingWindow) Size() int {
	return sw.windowSize
}

func (sw *SlidingWindow) Length() int {
	return len(sw.dataPoints)
}

// Add adds a new datapoint to the sliding window, if it's not full. Otherwise, it removes the oldest datapoint.
// Whenever it adds or removes a datapoint, it updates the total volume and total price by calling
// VolumeWeightedAveragePriceCalculator calculator's add and remove methods, and the result is stored in the
// attached calculator.
func (sw *SlidingWindow) Add(dataPoint DataPoint) {
	mux.Lock()
	defer mux.Unlock()

	if len(sw.dataPoints) == sw.windowSize {
		front := sw.dataPoints[0]
		back := sw.dataPoints[1:]

		sw.calculator.RemoveVolumeWeightedPrice(front.Price, front.Size)
		sw.dataPoints = back
	}

	sw.dataPoints = append(sw.dataPoints, dataPoint)
	sw.calculator.AddVolumeWeightedPrice(dataPoint.Price, dataPoint.Size)
}

func (sw *SlidingWindow) GetCalculator() *utils.VolumeWeightedAveragePriceCalculator {
	return &sw.calculator
}
