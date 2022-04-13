package vwap

import (
	"math/big"
	"sync"

	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/vwap/utils"
)

type SlidingWindow struct {
	mux          sync.Mutex
	currencyPair string
	dataPoints   []DataPoint
	windowSize   int
	calculator   utils.VolumeWeightedAveragePriceCalculator
}

type DataPoint struct {
	Type      string
	Size      *big.Float
	Price     *big.Float
	ProductId string
}

func NewSlidingWindow(maxSize int, currencyPair string) *SlidingWindow {
	return &SlidingWindow{
		currencyPair: currencyPair,
		dataPoints:   []DataPoint{},
		windowSize:   maxSize,
		calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
	}
}

func (sw *SlidingWindow) SetSize(maxSize int) {
	sw.windowSize = maxSize
}

func (sw *SlidingWindow) GetSize() int {
	return sw.windowSize
}

func (sw *SlidingWindow) GetLength() int {
	return len(sw.dataPoints)
}

func (sw *SlidingWindow) Add(dataPoint DataPoint) {
	sw.mux.Lock()
	defer sw.mux.Unlock()

	if len(sw.dataPoints) == sw.windowSize {
		front := sw.dataPoints[0]
		back := sw.dataPoints[1:]

		var transport []DataPoint

		transport = append(transport, back...)
		sw.dataPoints = transport

		sw.calculator.RemoveVolumeWeightedPrice(front.Price, front.Size)
	}

	sw.dataPoints = append(sw.dataPoints, dataPoint)
	sw.calculator.AddVolumeWeightedPrice(dataPoint.Price, dataPoint.Size)
}

func (sw *SlidingWindow) GetCalculator() *utils.VolumeWeightedAveragePriceCalculator {
	return &sw.calculator
}
