//go:build all
// +build all

package utils

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNewVolumeWeightedAveragePriceCalculator(t *testing.T) {
	tests := []struct {
		name string
		want *VolumeWeightedAveragePriceCalculator
	}{
		// Add TestNewVolumeWeightedAveragePriceCalculator test cases.
		{
			name: "Test New Volume Weighted Average Price Calculator with default values",
			want: &VolumeWeightedAveragePriceCalculator{
				ValueSum:  big.NewFloat(0),
				VolumeSum: big.NewFloat(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVolumeWeightedAveragePriceCalculator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVolumeWeightedAveragePriceCalculator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVolumeWeightedAveragePriceCalculator_AddVolumeWeightedPrice(t *testing.T) {
	type fields struct {
		valueSum  *big.Float
		volumeSum *big.Float
	}
	type args struct {
		price  *big.Float
		volume *big.Float
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// Add TestVolumeWeightedAveragePriceCalculator_AddVolumeWeightedPrice test cases.
		{
			name: "Test Volume Weighted Average Price Calculator Add Volume Weighted Price with default values",
			fields: fields{
				valueSum:  big.NewFloat(0),
				volumeSum: big.NewFloat(0),
			},
			args: args{
				price:  big.NewFloat(0),
				volume: big.NewFloat(0),
			},
		},
		{
			name: "Test Volume Weighted Average Price Calculator Add Volume Weighted Price with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(0),
				volumeSum: big.NewFloat(0),
			},
			args: args{
				price:  big.NewFloat(1),
				volume: big.NewFloat(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VolumeWeightedAveragePriceCalculator{
				ValueSum:  tt.fields.valueSum,
				VolumeSum: tt.fields.volumeSum,
			}
			v.AddVolumeWeightedPrice(tt.args.price, tt.args.volume)
		})
	}
}

func TestVolumeWeightedAveragePriceCalculator_Avg(t *testing.T) {
	type fields struct {
		valueSum  *big.Float
		volumeSum *big.Float
	}
	tests := []struct {
		name   string
		fields fields
		want   *big.Float
	}{
		// Add TestVolumeWeightedAveragePriceCalculator_Avg test cases.
		{
			name: "Test Volume Weighted Average Price Calculator Avg with default values",
			fields: fields{
				valueSum:  big.NewFloat(0),
				volumeSum: big.NewFloat(0),
			},
			want: big.NewFloat(0),
		},
		{
			name: "Test Volume Weighted Average Price Calculator Avg with zero volume",
			fields: fields{
				valueSum:  big.NewFloat(25),
				volumeSum: big.NewFloat(0),
			},
			want: big.NewFloat(0),
		},
		{
			name: "Test Volume Weighted Average Price Calculator Avg with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(1),
				volumeSum: big.NewFloat(1),
			},
			want: big.NewFloat(1),
		},
		{
			name: "Test Volume Weighted Average Price Calculator Avg with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(2),
				volumeSum: big.NewFloat(2),
			},
			want: big.NewFloat(1),
		},
		{
			name: "Test Volume Weighted Average Price Calculator Avg with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(3),
				volumeSum: big.NewFloat(3),
			},
			want: big.NewFloat(1),
		},
		{
			name: "Test Volume Weighted Average Price Calculator Avg with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(4),
				volumeSum: big.NewFloat(4),
			},
			want: big.NewFloat(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VolumeWeightedAveragePriceCalculator{
				ValueSum:  tt.fields.valueSum,
				VolumeSum: tt.fields.volumeSum,
			}
			if got := v.Avg(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVolumeWeightedAveragePriceCalculator_RemoveVolumeWeightedPrice(t *testing.T) {
	type fields struct {
		valueSum  *big.Float
		volumeSum *big.Float
	}
	type args struct {
		price  *big.Float
		volume *big.Float
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// Add TestVolumeWeightedAveragePriceCalculator_RemoveVolumeWeightedPrice cases with big.Float random values.
		{
			name: "Test Volume Weighted Average Price Calculator Remove Volume Weighted Price with default values",
			fields: fields{
				valueSum:  big.NewFloat(0),
				volumeSum: big.NewFloat(0),
			},
			args: args{
				price:  big.NewFloat(0),
				volume: big.NewFloat(0),
			},
		},
		{
			name: "Test Volume Weighted Average Price Calculator Remove Volume Weighted Price with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(25),
				volumeSum: big.NewFloat(15),
			},
			args: args{
				price:  big.NewFloat(1),
				volume: big.NewFloat(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VolumeWeightedAveragePriceCalculator{
				ValueSum:  tt.fields.valueSum,
				VolumeSum: tt.fields.volumeSum,
			}
			v.RemoveVolumeWeightedPrice(tt.args.price, tt.args.volume)
		})
	}
}

func TestVolumeWeightedAveragePriceCalculator_calculateWeightedAveragePrice(t *testing.T) {
	type fields struct {
		valueSum  *big.Float
		volumeSum *big.Float
	}
	tests := []struct {
		name   string
		fields fields
		want   *big.Float
	}{
		// Add TestVolumeWeightedAveragePriceCalculator_calculateWeightedAveragePrice test cases with big.Float random values.
		{
			name: "Test Volume Weighted Average Price Calculator calculate Weighted Average Price with default values",
			fields: fields{
				valueSum:  big.NewFloat(0),
				volumeSum: big.NewFloat(0),
			},
			want: big.NewFloat(0),
		},
		{
			name: "Test Volume Weighted Average Price Calculator calculate Weighted Average Price with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(25.1345),
				volumeSum: big.NewFloat(15.5768),
			},
			want: big.NewFloat(1.6135855888244055),
		},
		{
			name: "Test Volume Weighted Average Price Calculator calculate Weighted Average Price with non-zero values",
			fields: fields{
				valueSum:  big.NewFloat(0.5345),
				volumeSum: big.NewFloat(0.5768),
			},
			want: big.NewFloat(0.9266643550624133),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VolumeWeightedAveragePriceCalculator{
				ValueSum:  tt.fields.valueSum,
				VolumeSum: tt.fields.volumeSum,
			}
			if got := v.calculateWeightedAveragePrice(); !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("calculateWeightedAveragePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
