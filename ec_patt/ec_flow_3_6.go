package ec_patt

import (
	"scanner/schemas"
)

func GetEc(data schemas.AllSymData, indicatorData schemas.AllIndicatorsData) (map[string]string, map[string]string) {

	long_ec := make(map[string]string)
	short_ec := make(map[string]string)

	//loop through symbols
	for symbol, singleSymData := range data {

		//get req data
		macd_3_6 := indicatorData.Macd_3_6_9.Macd[symbol]
		signal_3_6_9 := indicatorData.Macd_3_6_9.Signal[symbol]
		don_5_up := indicatorData.DonchianChannel_5.Upper[symbol]
		don_5_low := indicatorData.DonchianChannel_5.Lower[symbol]
		ema_18 := indicatorData.Ema_18[symbol]

		long_ec_type := ec_flow(true, singleSymData, macd_3_6, signal_3_6_9, don_5_up, don_5_low, ema_18)
		short_ec_type := ec_flow(false, singleSymData, macd_3_6, signal_3_6_9, don_5_up, don_5_low, ema_18)

		if long_ec_type != "" {
			long_ec[symbol] = long_ec_type
		}
		if short_ec_type != "" {
			short_ec[symbol] = short_ec_type
		}

	}

	return long_ec, short_ec
}

func ec_flow(long bool, singleSymData schemas.SingleSymData, macd_3_6, signal_3_6_9, don_5_up, don_5_low, ema_18 []float64) string {

	findExtraDonBoIndex := func(long bool, data schemas.SingleSymData, lastDonBoIndex int, don_5_up, don_5_low []float64) int {
		data_len := len(data)
		//if we have bo after final ec
		for i := lastDonBoIndex + 1; i < data_len; i++ {
			bo_condition := data[i].Close > don_5_up[i]
			if !long {
				bo_condition = data[i].Close < don_5_low[i]
			}

			if bo_condition {
				return i
			}
		}
		return -1
	}

	data_len := len(singleSymData)

	trend_present := macd_3_6[data_len-1] > 0 || signal_3_6_9[data_len-1] > 0
	if !long {
		trend_present = macd_3_6[data_len-1] < 0 || signal_3_6_9[data_len-1] < 0
	}

	//trend present
	if !trend_present {
		return ""
	}

	//--pre ec--
	ec_type, firstDonBoIndex, lastDonBoIndex := pre_ec(long, singleSymData, macd_3_6, signal_3_6_9, don_5_up, don_5_low, ema_18)

	//no ec no need to continue
	if ec_type == "" {
		return ""
		//if ema is recent bar no need to continue for other ec
	} else if ec_type == "ema" && firstDonBoIndex == data_len-1 {
		return "pre_ec"
	}

	//--ec--
	firstDonBoIndex, lastDonBoIndex = ec(long, singleSymData, signal_3_6_9, don_5_up, don_5_low, ec_type, firstDonBoIndex, lastDonBoIndex)

	// if we don't have bo, so we have pre_ec only
	if lastDonBoIndex == -1 || firstDonBoIndex == -1 {
		return "pre_ec"
		//if last bo is recent bar no need to continure for other ec
	} else if lastDonBoIndex == data_len-1 {
		return "ec"
	}

	//--final ec---
	firstDonBoIndex, lastDonBoIndex = final_ec(long, singleSymData, don_5_up, don_5_low, lastDonBoIndex)
	//if we don't have bo, so we have pre_ec only
	if lastDonBoIndex == -1 || firstDonBoIndex == -1 {
		return "ec"
	} else if lastDonBoIndex == data_len-1 {
		return "final_ec"
	}

	//don bo after final ec
	extraDonBoIndex := findExtraDonBoIndex(long, singleSymData, lastDonBoIndex, don_5_up, don_5_low)

	//don't have extra bo after ec
	if extraDonBoIndex == -1 {
		return "final_ec"
	}

	return ""
}

func pre_ec(long bool, singleSymData schemas.SingleSymData, macd_3_6, signal_3_6, don_5_up, don_5_low, ema_18 []float64) (string, int, int) {
	data_len := len(singleSymData)

	macd_trend := false
	signal_trend := false
	if long {
		macd_trend = macd_3_6[data_len-1] > 0
		signal_trend = signal_3_6[data_len-1] > 0
	} else {
		macd_trend = macd_3_6[data_len-1] < 0
		signal_trend = signal_3_6[data_len-1] < 0
	}

	macd_bo_index := -1
	//if macd 3-6 trend , no signal trend
	if macd_trend && !signal_trend {
		macd_bo_index = findLinesBoIndex(long, macd_3_6, data_len-1)
		//singal trend (no care about macd trend)
	} else if signal_trend {
		signal_bo_index := findLinesBoIndex(long, signal_3_6, data_len-1)
		if signal_bo_index == -1 {
			return "", -1, -1
		}
		macd_bo_index = findLinesBoIndex(long, macd_3_6, signal_bo_index)
	} else {
		return "", -1, -1
	}

	if macd_bo_index == -1 {
		return "", -1, -1
	}

	//find don bo indexes
	firstDonBoIndex, lastDonBoIndex := findDonBoIndexes(long, don_5_up, don_5_low, singleSymData, macd_bo_index)

	if firstDonBoIndex != -1 && lastDonBoIndex != -1 {
		return "don", firstDonBoIndex, lastDonBoIndex
	}

	//no don bo : find ema 18 bo
	emaBoIndex := findEmaBO(long, singleSymData, ema_18, macd_bo_index)

	return "ema", emaBoIndex, -1
}

func ec(long bool, singleSymData schemas.SingleSymData, signal_3_6_9, don_5_up, don_5_low []float64, ec_type string, PrevFirstDonBoIndex, PrevLastDonBoIndex int) (int, int) {

	data_len := len(singleSymData)
	//between first don bo and last don bo (included) , we have signal bo >> it is ec not pre ec
	//so we can return same first and last don bo
	if ec_type == "don" {

		//signal bo before don bo
		trend_before_don_bo := signal_3_6_9[PrevFirstDonBoIndex] > 0
		if !long {
			trend_before_don_bo = signal_3_6_9[PrevFirstDonBoIndex] < 0
		}

		if trend_before_don_bo {
			return PrevFirstDonBoIndex, PrevLastDonBoIndex
		}

		//signal bo in impulse
		for i := PrevFirstDonBoIndex; i <= PrevLastDonBoIndex; i++ {
			signal_bo := signal_3_6_9[i] > 0 && signal_3_6_9[i-1] <= 0
			if !long {
				signal_bo = signal_3_6_9[i] < 0 && signal_3_6_9[i-1] >= 0
			}

			if signal_bo {
				return PrevFirstDonBoIndex, PrevLastDonBoIndex
			}
		}
	}

	//signal trend present
	if signal_3_6_9[len(singleSymData)-1] < 0 {
		return -1, -1
	}

	signal_bo_index := findLinesBoIndex(long, signal_3_6_9, data_len-1)

	if signal_bo_index == -1 {
		return -1, -1
	}

	//find don bo indexes
	firstDonBoIndex, lastDonBoIndex := findDonBoIndexes(long, don_5_up, don_5_low, singleSymData, signal_bo_index)

	return firstDonBoIndex, lastDonBoIndex

}

func final_ec(long bool, singleSymData schemas.SingleSymData, don_5_up, don_5_low []float64, PrevLastDonBoIndex int) (int, int) {

	//find don bo indexes
	return findDonBoIndexes(long, don_5_up, don_5_low, singleSymData, PrevLastDonBoIndex+1)
}

// helper functions
func findLinesBoIndex(long bool, macd_3_6 []float64, start int) int {
	for i := start; i > 0; i-- {

		bo_condition := macd_3_6[i] > 0 && macd_3_6[i-1] <= 0
		if !long {
			bo_condition = macd_3_6[i] < 0 && macd_3_6[i-1] >= 0
		}

		if bo_condition {
			return i
		}
	}
	return -1
}

func findEmaBO(long bool, data schemas.SingleSymData, ema_18 []float64, start int) int {

	for i := start; i < len(data); i++ {
		ema_bo := data[i].Close > ema_18[i]
		if !long {
			ema_bo = data[i].Close < ema_18[i]
		}

		if ema_bo {
			return i
		}
	}
	return -1
}

func findDonBoIndexes(long bool, don_up, don_low []float64, data schemas.SingleSymData, start int) (int, int) {

	firstDonBoIndex := -1
	lastDonBoIndex := -1

	findNewDonBoIndex := func(start int) int {
		for i := start; i < len(data); i++ {
			donBo := data[i].Close > don_up[i]
			if !long {
				donBo = data[i].Close < don_low[i]
			}

			if donBo {
				return i
			}
		}
		return -1
	}

	findLastDonBoIndex := func(firstDonBoIndex int) int {
		for i := firstDonBoIndex + 1; i < len(data); i++ {
			noDonBo := data[i].Close <= don_up[i]
			if !long {
				noDonBo = data[i].Close > don_low[i]
			}

			if noDonBo {
				return i - 1
			}
		}
		return len(data) - 1
	}

	//find first don bo
	firstDonBoIndex = findNewDonBoIndex(start)

	if firstDonBoIndex == -1 {
		return -1, -1
	}

	//find last don bo
	lastDonBoIndex = findLastDonBoIndex(firstDonBoIndex)

	//refine: extra buffer if
	// we don't have structral reaction between last don bo and new don bo
	//we can change last don bo
	if lastDonBoIndex < len(data)-1 {
		newDonBoIndex := findNewDonBoIndex(lastDonBoIndex + 1)

		if newDonBoIndex == -1 {
			return firstDonBoIndex, lastDonBoIndex
		}

		//don't want structral reaction between last don bo and new don bo
		for i := lastDonBoIndex + 1; i < newDonBoIndex; i++ {
			structralReaction := data[i].High < data[i-1].High
			if !long {
				structralReaction = data[i].Low > data[i-1].Low
			}

			if structralReaction {
				return firstDonBoIndex, lastDonBoIndex
			}
		}

		//new last don bo
		lastDonBoIndex = findLastDonBoIndex(newDonBoIndex)

	}

	return firstDonBoIndex, lastDonBoIndex
}
