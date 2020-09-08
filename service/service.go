package service

import (
	"currency/database"
	"sync"
)

const (
	//RegexApiTrade /api/trade/USD/GBP/5/2
	RegexApiTrade = `/api/trade/[a-zA-Z]{3}/[a-zA-Z]{3}/[0-9]+/[0-9]+`
	//RegexApiRegister /api/register/nick_name/email@qwe.com
	RegexApiRegister = `/api/register/[a-zA-Z0-9_\.]+/[a-zA-Z0-9_\.]+@[a-z]+\.com`
	//RegexApiInform /api/inform/nick_name/USD/GBP/5/2
	RegexApiInform   = `/api/inform/[a-zA-Z0-9_\.]+/[a-zA-Z]{3}/[a-zA-Z]{3}/[0-9]+/[0-9]+`
	currencyBasePath = "currency"
)

type MedianAnswer struct {
	Trade            string  `json:"trade"`
	MedianDifference float32 `json:"median_difference"`
}

type MovingAverageAnswer struct {
	state string
	mux   sync.Mutex
}

type Connection struct {
	database.Connection
}
