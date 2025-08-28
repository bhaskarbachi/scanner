package don_patt

import (
	"scanner/schemas"
)

func Recent_don5_bo(weeklyData schemas.AllSymData, weeklyIndicators schemas.AllIndicatorsData) ([]string, []string) {

	long_symbols := []string{}
	short_symbols := []string{}
	//loop through weekly data
	for symbol, singleSymData := range weeklyData {

		//get req data
		don_5_up := weeklyIndicators.DonchianChannel_5.Upper[symbol]
		don_5_low := weeklyIndicators.DonchianChannel_5.Lower[symbol]

		//check for recent donchian breakouts
		if long_don_bo(singleSymData, don_5_up) {
			long_symbols = append(long_symbols, symbol)
		}
		if short_don_bo(singleSymData, don_5_low) {
			short_symbols = append(short_symbols, symbol)
		}

	}

	return long_symbols, short_symbols
}

func long_don_bo(singleSymData schemas.SingleSymData, don_5_up []float64) bool {
	data_len := len(singleSymData)
	//check for recent donchian breakouts
	boType1 := singleSymData[data_len-1].Close > don_5_up[data_len-1]
	boType2 := singleSymData[data_len-2].Close > don_5_up[data_len-2] && singleSymData[data_len-1].Close < don_5_up[data_len-1]
	return boType1 || boType2
}

func short_don_bo(singleSymData schemas.SingleSymData, don_5_low []float64) bool {
	data_len := len(singleSymData)
	//check for recent donchian breakouts
	boType1 := singleSymData[data_len-1].Close < don_5_low[data_len-1]
	boType2 := singleSymData[data_len-2].Close < don_5_low[data_len-2] && singleSymData[data_len-1].Close > don_5_low[data_len-1]
	return boType1 || boType2
}
