package main

import (
	"scanner/all_scans"
	"scanner/data_process"
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

	//send data to all scans
	all_scans.ControlScans(dailyData, weeklyData, dailyIndicators, weeklyIndicators)

}
