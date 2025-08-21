package indicators

import (
	"scanner/schemas"
)

// function to calculate donchian channel for all symbols in data
func getDonchianChannel(data schemas.AllSymData, nbars int) (schemas.DonchianChannelData, error) {

	//initialize return data
	upper := make(map[string][]float64)
	lower := make(map[string][]float64)

	for symbol, singleSymData := range data {

		tempUpper, tempLower, err := donchianChannel(singleSymData, nbars)
		if err != nil {
			return schemas.DonchianChannelData{}, err
		}
		upper[symbol] = tempUpper
		lower[symbol] = tempLower

	}
	return schemas.DonchianChannelData{
		Upper: upper,
		Lower: lower,
	}, nil
}

func donchianChannel(data schemas.SingleSymData, nbars int) ([]float64, []float64, error) {
	dataLength := len(data)

	if dataLength < nbars {
		return []float64{}, []float64{}, nil
	}

	barNupper := make([]float64, dataLength)
	barNlower := make([]float64, dataLength)

	// Calculate initial values for the first nbars elements
	max, min := data[0].High, data[0].Low
	for i := 0; i < nbars; i++ {
		if data[i].High > max {
			max = data[i].High
		}
		if data[i].Low < min {
			min = data[i].Low
		}
		barNupper[i], barNlower[i] = max, min
	}

	// Calculate Donchian Channel values for the remaining elements
	for i := nbars; i < dataLength; i++ {
		barNupper[i], barNlower[i] = findMaxMinInRange(data, i-nbars, i)
	}

	return barNupper, barNlower, nil
}

func findMaxMinInRange(data schemas.SingleSymData, start, end int) (float64, float64) {
	max, min := data[start].High, data[start].Low
	for i := start + 1; i < end; i++ {
		if data[i].High > max {
			max = data[i].High
		}
		if data[i].Low < min {
			min = data[i].Low
		}
	}
	return max, min
}
