package vwap

import (
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"

	"bitbucket.org/keynear/coinbase-vwap-calculation/internal/vwap/utils"
)

func TestNewSlidingWindow(t *testing.T) {
	type args struct {
		maxSize      int
		currencyPair string
	}
	tests := []struct {
		name string
		args args
		want *SlidingWindow
	}{
		// Add TestNewSlidingWindow test cases.
		{
			name: "TestNewSlidingWindow",
			args: args{
				maxSize:      10,
				currencyPair: "BTC-USD",
			},
			want: &SlidingWindow{
				currencyPair: "BTC-USD",
				dataPoints:   []DataPoint{},
				windowSize:   10,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlidingWindow(tt.args.maxSize, tt.args.currencyPair); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewSlidingWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_Add(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	type args struct {
		dataPoint DataPoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// Add TestSlidingWindow_Add test cases with random data.
		{
			name: "TestSlidingWindow_Add",
			fields: fields{
				currencyPair: "BTC-USD",
				dataPoints:   []DataPoint{},
				windowSize:   10,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
			args: args{
				dataPoint: DataPoint{
					Price:     big.NewFloat(100.0),
					Size:      big.NewFloat(1.355559),
					Type:      "match",
					ProductID: "BTC-USD",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}
			sw.Add(tt.args.dataPoint)
		})
	}
}

func TestSlidingWindow_Add_Multiple(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	type args struct {
		dataPoints []DataPoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// Add TestSlidingWindow_Add_Multiple test cases with random data.
		{
			name: "TestSlidingWindow_Add_Multiple",
			fields: fields{
				currencyPair: "BTC-USD",
				dataPoints:   []DataPoint{},
				windowSize:   5,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
			args: args{
				dataPoints: []DataPoint{
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(1.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
				},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}

			for _, dataPoint := range tt.args.dataPoints {
				sw.Add(dataPoint)
			}
			if got := sw.Length(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_GetCalculator(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	tests := []struct {
		name   string
		fields fields
		want   *utils.VolumeWeightedAveragePriceCalculator
	}{
		// Add TestSlidingWindow_GetCalculator test cases.
		{
			name: "TestSlidingWindow_GetCalculator",
			fields: fields{
				currencyPair: "BTC-USD",
				dataPoints:   []DataPoint{},
				windowSize:   10,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
			want: &utils.VolumeWeightedAveragePriceCalculator{
				ValueSum:  big.NewFloat(0),
				VolumeSum: big.NewFloat(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}
			if got := sw.GetCalculator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCalculator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_Add_Avg(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	type args struct {
		dataPoints []DataPoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *big.Float
	}{
		// Add TestSlidingWindow_Add_Multiple test cases with random data.
		{
			name: "TestSlidingWindow_Add_Avg",
			fields: fields{
				currencyPair: "BTC-USD",
				dataPoints:   []DataPoint{},
				windowSize:   5,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
			args: args{
				dataPoints: []DataPoint{
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(0.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(3100.0),
						Size:      big.NewFloat(2.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(150.0),
						Size:      big.NewFloat(12.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(5100.0),
						Size:      big.NewFloat(25.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(6100.0),
						Size:      big.NewFloat(2.352239),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(5100.0),
						Size:      big.NewFloat(21.661324),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(31100.0),
						Size:      big.NewFloat(21.32559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(51100.0),
						Size:      big.NewFloat(41.255559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(61100.0),
						Size:      big.NewFloat(71.64587),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(81100.0),
						Size:      big.NewFloat(92.25539),
						Type:      "match",
						ProductID: "BTC-USD",
					},
				},
			},
			want: big.NewFloat(59406.42657),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}

			for _, dataPoint := range tt.args.dataPoints {
				sw.Add(dataPoint)
			}

			if got := sw.GetCalculator().Avg().String(); !reflect.DeepEqual(got, tt.want.String()) {
				t.Errorf("GetCalculator().Avg().String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_Add_Concurrent_Avg(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	type args struct {
		dataPoints []DataPoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *big.Float
	}{
		// Add TestSlidingWindow_Add_Multiple test cases with random data.
		{
			name: "TestSlidingWindow_Add_Concurrent_Avg",
			fields: fields{
				currencyPair: "BTC-USD",
				dataPoints:   []DataPoint{},
				windowSize:   5,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
			args: args{
				dataPoints: []DataPoint{
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(0.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(3100.0),
						Size:      big.NewFloat(2.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(150.0),
						Size:      big.NewFloat(12.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(5100.0),
						Size:      big.NewFloat(25.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(6100.0),
						Size:      big.NewFloat(2.352239),
						Type:      "match",
						ProductID: "BTC-USD",
					},
				},
			},
			want: big.NewFloat(3573.465985),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}

			count := len(tt.args.dataPoints)
			var wg sync.WaitGroup
			wg.Add(count)
			for _, dataPoint := range tt.args.dataPoints {
				go func(dataPoint DataPoint) {
					sw.Add(dataPoint)
					wg.Done()
				}(dataPoint)
			}
			wg.Wait()

			if got := sw.Length(); !reflect.DeepEqual(got, tt.fields.windowSize) {
				t.Errorf("Length() = %v, want %v", got, tt.fields.windowSize)
			}

			if got := sw.GetCalculator().Avg().String(); !reflect.DeepEqual(got, tt.want.String()) {
				t.Errorf("GetCalculator().Avg().String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_Add_Concurrent_Size(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	type args struct {
		dataPoints []DataPoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// Add TestSlidingWindow_Add_Multiple test cases with random data.
		{
			name: "TestSlidingWindow_Add_Concurrent_Size",
			fields: fields{
				currencyPair: "BTC-USD",
				dataPoints:   make([]DataPoint, 0),
				windowSize:   2,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
			args: args{
				dataPoints: []DataPoint{
					{
						Price:     big.NewFloat(100.0),
						Size:      big.NewFloat(0.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(3100.0),
						Size:      big.NewFloat(2.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(150.0),
						Size:      big.NewFloat(12.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(5100.0),
						Size:      big.NewFloat(25.355559),
						Type:      "match",
						ProductID: "BTC-USD",
					},
					{
						Price:     big.NewFloat(6100.0),
						Size:      big.NewFloat(2.352239),
						Type:      "match",
						ProductID: "BTC-USD",
					},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}

			var wg sync.WaitGroup
			testLength := len(tt.args.dataPoints)
			fmt.Println("testLength", testLength)
			wg.Add(testLength)

			for _, dataPoint := range tt.args.dataPoints {
				dataPoint := dataPoint
				go func() {
					fmt.Println("dataPoint", dataPoint)
					sw.Add(dataPoint)
					defer wg.Done()
				}()
			}

			wg.Wait()

			if got := sw.Length(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_SetSize(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	type args struct {
		maxSize int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// Add TestSlidingWindow_SetSize test cases.
		{
			name: "TestSlidingWindow_SetSize",
			fields: fields{
				currencyPair: "BTC-USD",
				dataPoints:   []DataPoint{},
				windowSize:   10,
				calculator:   *utils.NewVolumeWeightedAveragePriceCalculator(),
			},
			args: args{
				maxSize: 20,
			},
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}
			sw.SetSize(tt.args.maxSize)
			if got := sw.Size(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCalculator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_GetSize(t *testing.T) {
	type fields struct {
		currencyPair string
		dataPoints   []DataPoint
		windowSize   int
		calculator   utils.VolumeWeightedAveragePriceCalculator
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// Add TestSlidingWindow_GetSize test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SlidingWindow{
				currencyPair: tt.fields.currencyPair,
				dataPoints:   tt.fields.dataPoints,
				windowSize:   tt.fields.windowSize,
				calculator:   tt.fields.calculator,
			}
			if got := sw.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
