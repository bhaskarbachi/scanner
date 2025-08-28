package don_patt

import "scanner/schemas"

func Recent_don5_cb(weeklyData schemas.AllSymData, weeklyIndicators schemas.AllIndicatorsData) ([]string, []string) {

	long_symbols := []string{}
	short_symbols := []string{}
	//loop through weekly data
	for symbol, singleSymData := range weeklyData {

		//get req data
		don_5_up := weeklyIndicators.DonchianChannel_5.Upper[symbol]
		don_5_low := weeklyIndicators.DonchianChannel_5.Lower[symbol]

		//get recent don bo index(bar 5 bo bar) and direction
		recentBoIndex, direction := getRecentDon5Breakout(don_5_up, don_5_low, singleSymData)

		//bo recent breakout or bo at recent bar
		if recentBoIndex == -1 || recentBoIndex == len(singleSymData)-1 {
			continue
		}

		if direction == "down" {
			if donCb(true, singleSymData, recentBoIndex, don_5_up, don_5_low) {
				long_symbols = append(long_symbols, symbol)
			}
		} else {
			if donCb(false, singleSymData, recentBoIndex, don_5_up, don_5_low) {
				short_symbols = append(short_symbols, symbol)
			}
		}

	}

	return long_symbols, short_symbols
}

func donCb(long bool, singleSymData schemas.SingleSymData, recentBoIndex int, don_5_up, don_5_low []float64) bool {

	data_len := len(singleSymData)

	if recentBoIndex == 0 {
		return false
	}

	isCb := func(level_val float64) bool {

		var cbType1, cbType1Equal, cbType2, cbType2Equal bool
		if long {

			cbType1 = singleSymData[data_len-1].Close > level_val && singleSymData[data_len-2].Close < level_val
			cbType1Equal = singleSymData[data_len-1].Close > level_val && singleSymData[data_len-2].Close == level_val && singleSymData[data_len-2].Open < level_val

			cbType2 = data_len > 2 && singleSymData[data_len-1].Close >= level_val && singleSymData[data_len-2].Close > level_val && singleSymData[data_len-3].Close < level_val
			cbType2Equal = data_len > 2 && singleSymData[data_len-1].Close >= level_val && singleSymData[data_len-2].Close > level_val && singleSymData[data_len-3].Close == level_val && singleSymData[data_len-3].Open < level_val

		} else {
			cbType1 = singleSymData[data_len-1].Close < level_val && singleSymData[data_len-2].Close > level_val
			cbType1Equal = singleSymData[data_len-1].Close < level_val && singleSymData[data_len-2].Close == level_val && singleSymData[data_len-2].Open > level_val

			cbType2 = data_len > 2 && singleSymData[data_len-1].Close <= level_val && singleSymData[data_len-2].Close < level_val && singleSymData[data_len-3].Close > level_val
			cbType2Equal = data_len > 2 && singleSymData[data_len-1].Close <= level_val && singleSymData[data_len-2].Close < level_val && singleSymData[data_len-3].Close == level_val && singleSymData[data_len-3].Open > level_val
		}
		return cbType1 || cbType1Equal || cbType2 || cbType2Equal
	}

	//find levels
	xlevel_val := findXLevel(long, singleSymData, recentBoIndex)
	first_don_bo_val := findFirstDonBoLevel(long, singleSymData, recentBoIndex, don_5_up, don_5_low)
	horz_bar_5_val := findHorzBar5Level(long, recentBoIndex, don_5_up, don_5_low)

	//cb on any level above
	if xlevel_val != -1 && isCb(xlevel_val) {
		return true
	} else if first_don_bo_val != -1 && isCb(first_don_bo_val) {
		return true
	} else if horz_bar_5_val != -1 && isCb(horz_bar_5_val) {
		return true
	}

	return false
}

//----helper functions----

// get recent don bo index and direction
func getRecentDon5Breakout(don_5_up, don_5_low []float64, singleSymData schemas.SingleSymData) (int, string) {

	for i := len(singleSymData) - 1; i >= 0; i-- {
		if singleSymData[i].Close > don_5_up[i] {
			return i, "up"
		}
		if singleSymData[i].Close < don_5_low[i] {
			return i, "down"
		}
	}
	return -1, ""
}

// find x level value before recent bo index
func findXLevel(long bool, singleSymData schemas.SingleSymData, recentBoIndex int) float64 {

	if long {
		return singleSymData[recentBoIndex-1].Low
	} else {
		return singleSymData[recentBoIndex-1].High
	}
}

// find don value where first breakout occurs
func findFirstDonBoLevel(long bool, singleSymData schemas.SingleSymData, recentBoIndex int, don_5_up, don_5_low []float64) float64 {

	if long {
		for i := recentBoIndex - 1; i >= 0; i-- {
			if singleSymData[i].Close >= don_5_low[i] {
				return don_5_low[i+1]
			}
		}
	} else {
		for i := recentBoIndex - 1; i >= 0; i-- {
			if singleSymData[i].Close <= don_5_up[i] {
				return don_5_up[i+1]
			}
		}

	}

	return -1
}

// find horizontal bar 5 level
func findHorzBar5Level(long bool, recentBoIndex int, don_5_up, don_5_low []float64) float64 {

	upto := recentBoIndex - 10
	if upto < 0 {
		upto = 0
	}

	if long {
		for i := recentBoIndex; i > upto; i-- {
			if don_5_low[i] == don_5_low[i-1] {
				return don_5_low[i]
			}
		}

	} else {
		for i := recentBoIndex; i > upto; i-- {
			if don_5_up[i] == don_5_up[i-1] {
				return don_5_up[i]
			}
		}
	}

	return -1
}
