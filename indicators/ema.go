package indicators

import (
	"fmt"
	"math"
	"scanner/schemas"
)

func getEma(data schemas.AllSymData, window int) (map[string][]float64, error) {

	//initialize return data
	emaValues := make(map[string][]float64)

	for symbol, singleSymData := range data {
		tempEma, err := emaGeneric(extractClosePrices(singleSymData), window)
		if err != nil {
			return nil, err
		}
		emaValues[symbol] = tempEma
	}

	return emaValues, nil
}

// Generic EMA function that works with any float64 slice
func emaGeneric(data []float64, window int) ([]float64, error) {
	if window <= 0 {
		return nil, fmt.Errorf("window size should be above 0, given window size: %v", window)
	}

	dataLength := len(data)
	if dataLength == 0 {
		return []float64{}, nil
	}

	emaValues := make([]float64, dataLength)
	alpha := 2.0 / float64(window+1)
	oneMinusAlpha := 1.0 - alpha

	// Fill initial values with NaN
	for i := 0; i < window-1 && i < dataLength; i++ {
		emaValues[i] = math.NaN()
	}

	// If we don't have enough data for even one EMA value
	if dataLength < window {
		return emaValues, nil
	}

	// Initialize with SMA
	emaValues[window-1] = average(data[:window])

	// Calculate EMA
	for i := window; i < dataLength; i++ {
		emaValues[i] = alpha*data[i] + oneMinusAlpha*emaValues[i-1]
	}

	return emaValues, nil
}

// Extract close prices from SingleSymData
func extractClosePrices(data schemas.SingleSymData) []float64 {
	prices := make([]float64, len(data))
	for i, candle := range data {
		prices[i] = candle.Close
	}
	return prices
}

// Calculate average of a float64 slice
func average(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}
