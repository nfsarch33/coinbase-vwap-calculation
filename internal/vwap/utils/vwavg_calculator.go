package utils

import (
	"math/big"
)

// VolumeWeightedAveragePriceCalculator is a struct that calculates the VWAP.
type VolumeWeightedAveragePriceCalculator struct {
	ValueSum  *big.Float
	VolumeSum *big.Float
}

func NewVolumeWeightedAveragePriceCalculator() *VolumeWeightedAveragePriceCalculator {
	return &VolumeWeightedAveragePriceCalculator{
		ValueSum:  big.NewFloat(0),
		VolumeSum: big.NewFloat(0),
	}
}

func NewBigFloat() *big.Float {
	return new(big.Float)
}

// Avg is the shorthand for calculateWeightedAveragePrice method.
func (v *VolumeWeightedAveragePriceCalculator) Avg() *big.Float {
	return v.calculateWeightedAveragePrice()
}

// calculateWeightedAveragePrice calculates the VWAP = (sum(value * volume) / sum(volume)).
func (v *VolumeWeightedAveragePriceCalculator) calculateWeightedAveragePrice() *big.Float {
	if v.VolumeSum.Cmp(big.NewFloat(0)) == 0 {
		return big.NewFloat(0)
	}

	return NewBigFloat().Quo(v.ValueSum, v.VolumeSum)
}

// AddVolumeWeightedPrice adds the value and volume to the calculator.
// ValuesSum = sum(value * volume) and VolumeSum = sum(volume).
func (v *VolumeWeightedAveragePriceCalculator) AddVolumeWeightedPrice(price *big.Float, volume *big.Float) {
	tempPrice := NewBigFloat()
	tempPrice.Add(v.ValueSum, tempPrice.Mul(price, volume))
	v.ValueSum = tempPrice

	tempVolume := NewBigFloat()
	tempVolume.Add(v.VolumeSum, volume)
	v.VolumeSum = tempVolume
}

// RemoveVolumeWeightedPrice removes the value and volume from the calculator.
// ValuesSum = sum(value * volume) - value * volume, and VolumeSum = sum(volume) - volume.
func (v *VolumeWeightedAveragePriceCalculator) RemoveVolumeWeightedPrice(price *big.Float, volume *big.Float) {
	tempPrice := NewBigFloat()
	tempPrice.Sub(v.ValueSum, tempPrice.Mul(price, volume))
	v.ValueSum = tempPrice

	tempVolume := NewBigFloat()
	tempVolume.Sub(v.VolumeSum, volume)
	v.VolumeSum = tempVolume
}
