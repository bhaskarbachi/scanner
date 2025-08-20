package schemas

import "time"

type TOHLCV struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

type SingleSymData []TOHLCV

type AllSymData map[string]SingleSymData
