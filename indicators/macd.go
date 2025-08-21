package indicators

import (
	"fmt"
	"math"
	"scanner/schemas"
)

func getMacd(data schemas.AllSymData, m1Length, m2Length, signalLength int) (schemas.MacdData, error) {
	// Initialize return data
	macdValues := make(map[string][]float64)
	signalValues := make(map[string][]float64)

	for symbol, singleSymData := range data {
		tempMacd, tempSignal, err := macd(singleSymData, m1Length, m2Length, signalLength)
		if err != nil {
			return schemas.MacdData{}, err
		}
		macdValues[symbol] = tempMacd
		signalValues[symbol] = tempSignal
	}

	return schemas.MacdData{
		Macd:   macdValues,
		Signal: signalValues,
	}, nil
}

func macd(data schemas.SingleSymData, m1Length, m2Length, signalLength int) ([]float64, []float64, error) {
	if m1Length >= m2Length {
		return nil, nil, fmt.Errorf("m1 length should be less than m2, but got m1: %v, m2: %v", m1Length, m2Length)
	}

	if len(data) == 0 {
		return []float64{}, []float64{}, nil
	}

	// Extract close prices
	closePrices := extractClosePrices(data)

	// Calculate EMAs
	m1Ema, err := emaGeneric(closePrices, m1Length)
	if err != nil {
		return nil, nil, fmt.Errorf("error during m1 EMA calculation: %w", err)
	}

	m2Ema, err := emaGeneric(closePrices, m2Length)
	if err != nil {
		return nil, nil, fmt.Errorf("error during m2 EMA calculation: %w", err)
	}

	// Calculate MACD line
	macdLine := make([]float64, len(data))
	for i := range macdLine {
		if math.IsNaN(m1Ema[i]) || math.IsNaN(m2Ema[i]) {
			macdLine[i] = math.NaN()
		} else {
			macdLine[i] = m1Ema[i] - m2Ema[i]
		}
	}

	// Calculate signal line - we know MACD starts being valid at m2Length-1
	signalLine := make([]float64, len(data))
	validMacdStart := m2Length - 1

	// Fill all initial values with NaN
	for i := 0; i < len(signalLine); i++ {
		signalLine[i] = math.NaN()
	}

	// Only calculate signal if we have enough valid MACD values
	if len(data) > validMacdStart && len(data)-validMacdStart >= signalLength {
		validMacdValues := macdLine[validMacdStart:]
		tempSignal, err := emaGeneric(validMacdValues, signalLength)
		if err != nil {
			return nil, nil, fmt.Errorf("error calculating signal line: %w", err)
		}

		// Copy the calculated signal values to the correct positions
		copy(signalLine[validMacdStart:], tempSignal)
	}

	return macdLine, signalLine, nil
}
