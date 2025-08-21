package indicators

import "scanner/schemas"

func GetAllIndicatorsData(data schemas.AllSymData) schemas.AllIndicatorsData {

	//cal donchain indicators
	donchianChannel_5, err := getDonchianChannel(data, 5)
	if err != nil {
		panic(err)
	}

	//calucalting ema indicators
	ema_3, err := getEma(data, 3)
	if err != nil {
		panic(err)
	}
	ema_6, err := getEma(data, 6)
	if err != nil {
		panic(err)
	}
	ema_9, err := getEma(data, 9)
	if err != nil {
		panic(err)
	}
	ema_18, err := getEma(data, 18)
	if err != nil {
		panic(err)
	}
	ema_50, err := getEma(data, 50)
	if err != nil {
		panic(err)
	}

	//cal macd indicators
	macd_3_6_9, err := getMacd(data, 3, 6, 9)
	if err != nil {
		panic(err)
	}

	macd_12_26_9, err := getMacd(data, 12, 26, 9)
	if err != nil {
		panic(err)
	}

	macd_18_50_9, err := getMacd(data, 18, 50, 9)
	if err != nil {
		panic(err)
	}

	macd_50_100_9, err := getMacd(data, 50, 100, 9)
	if err != nil {
		panic(err)
	}

	macd_100_200_9, err := getMacd(data, 100, 200, 9)
	if err != nil {
		panic(err)
	}

	return schemas.AllIndicatorsData{
		DonchianChannel_5: donchianChannel_5,
		Ema_3:             ema_3,
		Ema_6:             ema_6,
		Ema_9:             ema_9,
		Ema_18:            ema_18,
		Ema_50:            ema_50,
		Macd_3_6_9:        macd_3_6_9,
		Macd_12_26_9:      macd_12_26_9,
		Macd_18_50_9:      macd_18_50_9,
		Macd_50_100_9:     macd_50_100_9,
		Macd_100_200_9:    macd_100_200_9,
	}
}
