package currency

import (
	"encoding/json"
	"time"
)

type CurrencyConnection interface {
	SetupCurrency()
}

func (con CurrConnection) SetupCurrency() {

	checkError(createCurrencyTable(con))
	checkError(createRatesTable(con))

	checkError(setInitialValuesForCurrencyTable(con))
	checkError(setInitialValuesForRatesTable(con, initializingRatesForDays))

	go con.UpdateDaily()
}

func createCurrencyTable(con CurrConnection) error {
	querry := `CREATE TABLE currencySQL.currency(
	id int not null auto_increment unique,
	name varchar(255),
	PRIMARY KEY(id),
	UNIQUE(name)
	);`
	_, err := con.DBConnect.Exec(querry)
	return err
}

func createRatesTable(con CurrConnection) error {
	querry := `CREATE TABLE currencySQL.rates(
		id1 int not null,
		id2 int not null,
		rate float8,
		creation_date DATETIME,
		FOREIGN KEY(id1) REFERENCES currency(id),
		FOREIGN KEY(id2) REFERENCES currency(id)
		);`
	_, err := con.DBConnect.Exec(querry)
	return err
}

//SetInitialValuesForRatesTable is function that takes one argument
//and initalizes rates table in currencySQL database for `days` days back
//but does not take rates from current day
func setInitialValuesForRatesTable(con CurrConnection, days int) error {
	currentTime := time.Now()
	beginDate := currentTime.AddDate(0, 0, -days)
	beginDateString := extractDate(beginDate)
	endDateString := extractDate(currentTime)
	var err error
	for beginDateString != endDateString {
		err = con.InsertHistoricalRatio(beginDateString)
		if err != nil {
			return err
		}
		beginDate = beginDate.AddDate(0, 0, 1)
		beginDateString = extractDate(beginDate)
	}
	err = con.InsertHistoricalRatio(beginDateString)
	return nil
}

//SetInitialValuesForCurrencyTable should be called only once
//function is initalizing table currency in database currencySQL
//with currency symbols ("USD","GBP","EUR"...)
func setInitialValuesForCurrencyTable(con CurrConnection) error {
	symbolsPath := buildPath(apiSymbolsEndpoint)
	bodyBites, err := runClient(symbolsPath)
	if err != nil {
		return err
	}
	var dataResponse SymbolsResponse
	if err := json.Unmarshal(bodyBites, &dataResponse); err != nil {
		return err
	}

	currList := make([]Currency, 0)
	for abberivation := range dataResponse.Symbols {
		var curr Currency
		curr.Name = abberivation
		currList = append(currList, curr)
	}
	if _, err := con.InsertCurrencyList(currList); err != nil {
		return err
	}
	return nil
}
