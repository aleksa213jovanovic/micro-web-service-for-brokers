package currency

import (
	"currency/database"
	"log"
	"time"
)

//_ "github.com/go-sql-driver/mysql"
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var initializingRatesForDays = 7

const (
	accessKey          = "?access_key=f6bfc9703e66618a7c4edbca2c45bfc3"
	apiBasePath        = "http://data.fixer.io/api"
	apiSymbolsEndpoint = "symbols"
	apiLatestEndpoint  = "latest"

	//if dateFormat is changed, extractDate function must be changed too
	dateFormat = "2006-01-02 15:04:05"
)

//SymbolsResponse is storing json object
//from get data.fixer.io/api/symbols response
type SymbolsResponse struct {
	Succes  bool              `json:"success"`
	Symbols map[string]string `json:"symbols"`
}

//RatesResponse is storing json object
//from get data.fixer.io/api/latest response
type RatesResponse struct {
	BaseCurrency string             `json:"base"`
	Date         string             `json:"date"`
	Rates        map[string]float32 `json:"rates"`
}

//Currency is struct mapped into db table currencies
type Currency struct {
	ID   int64
	Name string `json:"name"`
}

//Value is struct mapped into db table values
type Value struct {
	ID1  int64
	ID2  int64
	Rate float32   `json:"rate"`
	Date time.Time `json:"date"`
}

type CurrConnection struct {
	database.Connection
}
