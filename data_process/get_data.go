package data_process

import (
	"fmt"
	"scanner/schemas"
)

func GetData(file_path string) (schemas.AllSymData, schemas.AllSymData, []string, error) {

	//getting daily data from csv file
	dailyData, err := getDailyData(file_path)
	if err != nil {
		fmt.Println("error getting daily data:", err)
		return nil, nil, nil, err
	}

	//converting daily data to weekly data
	weeklyData, err := getWeeklyData(dailyData)
	if err != nil {
		fmt.Println("error converting daily data to weekly data:", err)
		return nil, nil, nil, err
	}

	//Get all symbols from daily data
	symbols := getSymbols(dailyData)

	return dailyData, weeklyData, symbols, nil
}

func getSymbols(data schemas.AllSymData) []string {
	symbols := make([]string, 0)
	for symbol := range data {
		symbols = append(symbols, symbol)
	}
	return symbols
}
