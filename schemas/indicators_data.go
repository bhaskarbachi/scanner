package schemas

// struct to hold donchian channel data of multiple symbols
type DonchianChannelData struct {
	Upper map[string][]float64
	Lower map[string][]float64
}

type MacdData struct {
	Macd   map[string][]float64
	Signal map[string][]float64
}

// type Recent don bo data
type DonBoData struct {
	Recent          map[string]int
	RecentDirection map[string]string
}

// struct hold all indicators data
type AllIndicatorsData struct {
	DonchianChannel_5  DonchianChannelData
	DonchianChannel_10 DonchianChannelData
	DonchianChannel_20 DonchianChannelData
	DonchianChannel_50 DonchianChannelData
	Ema_3              map[string][]float64
	Ema_6              map[string][]float64
	Ema_9              map[string][]float64
	Ema_18             map[string][]float64
	Ema_50             map[string][]float64
	Macd_3_6_9         MacdData
	Macd_12_26_9       MacdData
	Macd_18_50_9       MacdData
	Macd_50_100_9      MacdData
	Macd_100_200_9     MacdData
}
