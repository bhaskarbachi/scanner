package data_process

import (
	"scanner/schemas"
	"time"
)

// function to convert daily data to weekly data
func getWeeklyData(data schemas.AllSymData) (schemas.AllSymData, error) {

	//initialize return data
	weeklyData := make(schemas.AllSymData)

	//convert to weekly data
	for symbol, dailyData := range data {
		singleWeekData, err := convertToWeekly(dailyData)
		if err != nil {
			return nil, err
		}
		weeklyData[symbol] = singleWeekData
	}
	return weeklyData, nil
}

// function to convert single symbol daily data to weekly data
func convertToWeekly(data schemas.SingleSymData) (schemas.SingleSymData, error) {

	var weeklyData schemas.SingleSymData
	var currentWeekData schemas.SingleSymData

	for _, dailyData := range data {
		if len(currentWeekData) == 0 || sameWeek(currentWeekData[0].Date.UTC(), dailyData.Date.UTC()) {
			// Append to the current week's data if it's the same week
			currentWeekData = append(currentWeekData, dailyData)
		} else {
			// Process and store the completed week before starting the next week
			weeklyData = append(weeklyData, convertDailyToWeekly(currentWeekData))
			currentWeekData = schemas.SingleSymData{dailyData}
		}
	}

	// Directly append the last processed week
	if len(currentWeekData) > 0 {
		weeklyData = append(weeklyData, convertDailyToWeekly(currentWeekData))
	}

	return weeklyData, nil
}

// function to check if two dates are in same week
func sameWeek(date1, date2 time.Time) bool {
	//getting mondays of both weeks
	monday1 := startOfWeek(date1)
	monday2 := startOfWeek(date2)
	//return bool for:both mondays are same or not
	return monday1.Equal(monday2)
}

// helper function to get "date" of starting day of week(monday)
func startOfWeek(date time.Time) time.Time {
	//gettin weekday(sun, mon ,..) converting weekday into integer
	// ex: mon=1, tue=2,....sun=0
	n := int(date.Weekday())
	//for sunday we get 0, but we want 7 to make it to current week
	if n == 0 {
		n = 7
	}
	//returning monday date
	return date.AddDate(0, 0, -n+1)
}

// function to convert daily data to weekly
func convertDailyToWeekly(currentWeeklyData schemas.SingleSymData) schemas.TOHLCV {
	weekStartDate := currentWeeklyData[0].Date
	openPrice := currentWeeklyData[0].Open
	closePrice := currentWeeklyData[len(currentWeeklyData)-1].Close
	highPrice, lowPrice, totalVolume := maxAndMin(currentWeeklyData)

	return schemas.TOHLCV{
		Date:   weekStartDate,
		Open:   openPrice,
		High:   highPrice,
		Low:    lowPrice,
		Close:  closePrice,
		Volume: totalVolume,
	}
}

// function to get max and minium in dailydata
func maxAndMin(s schemas.SingleSymData) (float64, float64, int) {

	var max, min, sum = s[0].High, s[0].Low, s[0].Volume

	for _, v := range s[1:] {
		if v.High > max {
			max = v.High
		}

		if v.Low < min {
			min = v.Low
		}

		sum += v.Volume
	}

	return max, min, sum

}
