package main

import (
	"fmt"
	"os"
	"scanner/data_process"
	"scanner/ec_patt"
	"scanner/filters"
	"scanner/indicators"
)

func main() {

	//get req data
	//inWeek := in_week()

	//getting daily data and weekly data and symbols
	dailyData, weeklyData, symbols, err := data_process.GetData("daily_data.csv")
	if err != nil {
		panic("error getting data:" + err.Error())
	}

	// Calculate all daily and weekly indicators
	dailyIndicators := indicators.GetAllIndicatorsData(dailyData)
	weeklyIndicators := indicators.GetAllIndicatorsData(weeklyData)

	_, _, _ = dailyIndicators, weeklyIndicators, symbols

	// get ec symbols
	long_ec, short_ec := ec_patt.GetEc(dailyData, dailyIndicators)
	long_ec_keys := extraSyminKeys(long_ec)
	short_ec_keys := extraSyminKeys(short_ec)

	//filter long and short ec with structral reaction
	fil_long_ec_sym := filters.StructuralReaction(true, dailyData, long_ec_keys)
	fil_short_ec_sym := filters.StructuralReaction(false, dailyData, short_ec_keys)

	fmt.Println("Filtered Long EC Symbols:", len(fil_long_ec_sym), fil_long_ec_sym)
	fmt.Println()
	fmt.Println("Filtered Short EC Symbols:", len(fil_short_ec_sym), fil_short_ec_sym)

	//save it to TLS file
	writeSymbolsToTLS(fil_long_ec_sym, "long_ec.tls")
	writeSymbolsToTLS(fil_short_ec_sym, "short_ec.tls")

}

func extraSyminKeys(myDict map[string]string) []string {
	var keys []string
	for key := range myDict {
		keys = append(keys, key)
	}
	return keys
}

// writeSymbolsToTLS writes only the symbol values to a TLS file, one per line.
func writeSymbolsToTLS(symbols []string, filename string) error {
	file, err := os.Create(filename) // Simplified file creation
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	for _, value := range symbols {
		_, err := file.WriteString(value + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}
