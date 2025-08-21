package filters

import "scanner/schemas"

func StructuralReaction(long bool, data schemas.AllSymData, symbols []string) []string {

	isStructuralReaction := func(symbol string) bool {
		data_len := len(data[symbol])
		reaction := data[symbol][data_len-1].High < data[symbol][data_len-2].High
		if !long {
			reaction = data[symbol][data_len-1].Low > data[symbol][data_len-2].Low
		}
		return reaction
	}

	var result []string
	for _, symbol := range symbols {
		if isStructuralReaction(symbol) {
			result = append(result, symbol)
		}
	}
	return result
}
