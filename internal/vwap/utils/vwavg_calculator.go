package utils

import (
	"math/big"
)

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

func (v *VolumeWeightedAveragePriceCalculator) Avg() *big.Float {
	return v.calculateWeightedAveragePrice()
}

func (v *VolumeWeightedAveragePriceCalculator) calculateWeightedAveragePrice() *big.Float {
	if v.VolumeSum.Cmp(big.NewFloat(0)) == 0 {
		return big.NewFloat(0)
	}

	return NewBigFloat().Quo(v.ValueSum, v.VolumeSum)
}

func (v *VolumeWeightedAveragePriceCalculator) AddVolumeWeightedPrice(price *big.Float, volume *big.Float) {
	tempPrice := NewBigFloat()
	tempPrice.Add(v.ValueSum, tempPrice.Mul(price, volume))
	v.ValueSum = tempPrice

	tempVolume := NewBigFloat()
	tempVolume.Add(v.VolumeSum, volume)
	v.VolumeSum = tempVolume
}

func (v *VolumeWeightedAveragePriceCalculator) RemoveVolumeWeightedPrice(price *big.Float, volume *big.Float) {
	tempPrice := NewBigFloat()
	tempPrice.Sub(v.ValueSum, tempPrice.Mul(price, volume))
	v.ValueSum = tempPrice

	tempVolume := NewBigFloat()
	tempVolume.Sub(v.VolumeSum, volume)
	v.VolumeSum = tempVolume
}
