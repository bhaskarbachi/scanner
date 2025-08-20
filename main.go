package main

import (
	"fmt"
	"scanner/data_process"
)

func main() {

	//get req data
	//inWeek := in_week()

	dailyData, weeklyData, symbols, err := data_process.GetData("daily_data.csv")
	if err != nil {
		panic("error getting data:" + err.Error())
	}

	_, _, _ = dailyData, weeklyData, symbols

	fmt.Println(dailyData[symbols[0]])

}
