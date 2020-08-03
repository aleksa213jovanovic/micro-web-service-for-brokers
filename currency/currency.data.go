package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type CurrencyRepository interface {
	InsertRatio(dataResponse RatesResponse) (int64, error)
	InsertCurrencyList(currList []Currency) (int64, error)
	GetAllCurrencySymbols() (map[string]int64, error)
	InsertHistoricalRatio(fromDate string) error
	InsertLatestRatio() error
	MedianControler(baseCurrency string, tradeCurrency string, longTermMedian int, shortTermMedian int) (float32, float32, error)
	UpdateDaily()
}

//insertRatio takse one argument of type RatesResponse
//(RatesResponse is a struct in which we extract json response from get request)
//and inserts all unorderd pairs of currency symbols and their corresponding ratios
//into currencySQL.rates table
//retur values are number of rows affected and error
func (con CurrConnection) InsertRatio(dataResponse RatesResponse) (int64, error) {
	delete(dataResponse.Rates, dataResponse.BaseCurrency)
	allCurrencies, err := con.GetAllCurrencySymbols()
	if err != nil {
		return 0, err
	}

	querry := fmt.Sprintf("INSERT INTO currencySQL.rates (id1,id2,rate,creation_date) VALUES")
	for len(dataResponse.Rates) >= 1 {
		baseCurrID := allCurrencies[dataResponse.BaseCurrency]
		var tradeCurr string
		var currRate float32
		for tradeCurr, currRate = range dataResponse.Rates {
			if tradeCurr != dataResponse.BaseCurrency {
				querryPostfix := fmt.Sprintf(" (%d,%d,%g, '%s 00:00:00'),", baseCurrID, allCurrencies[tradeCurr], currRate, dataResponse.Date)
				querry += querryPostfix
			} else {
				fmt.Printf("%s %s", tradeCurr, dataResponse.BaseCurrency)
			}
		}

		if len(dataResponse.Rates) >= 1 {

			dataResponse, err = convertResponseRatesToNewBase(dataResponse, tradeCurr)
			if err != nil {
				return 0, err
			}
			delete(dataResponse.Rates, tradeCurr)
		}
	}
	querry = querry[0 : len(querry)-1]
	querry += ";"
	result, err := con.DBConnect.Exec(querry)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int64(rowsAffected), nil
}

//insertCurrencyList takes a slice of Currency structs and inserts every
//symbol from it
//returns number of rows affected and errror
func (con CurrConnection) InsertCurrencyList(currList []Currency) (int64, error) {
	query := `INSERT INTO currencySQL.currency (name) VALUES`
	for _, curr := range currList {
		query += fmt.Sprintf("(\"%s\"),", curr.Name)
	}
	query = query[0 : len(query)-1]
	result, err := con.DBConnect.Exec(query)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

//getAllCurrencySymbols will take all symbols and their ids from currencySQL.currency table
//returning a map where key is symbol string ('EUR','RSD'...) and value is corresponding id (primary_key)
func (con CurrConnection) GetAllCurrencySymbols() (map[string]int64, error) {
	result, err := con.DBConnect.Query(`SELECT id, name FROM currencySQL.currency`)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	ret := make(map[string]int64, 0)
	for result.Next() {
		var curr Currency
		result.Scan(&curr.ID, &curr.Name)
		ret[curr.Name] = curr.ID
	}
	return ret, nil
}

//InsertHistoricalRatio takes one argument 'fromDate'
//and inserts all currency rates from 'fromDate' date
//into a currencySQL.rates table
func (con CurrConnection) InsertHistoricalRatio(fromDate string) error {
	datePath := buildPath(fromDate)

	clientResponse, err := runClient(datePath)
	if err != nil {
		return err
	}
	dataResponse, err := convertClientResponse(clientResponse)
	if err != nil {
		return err
	}
	if _, err := con.InsertRatio(dataResponse); err != nil {
		return err
	}
	return nil
}

var i = 1

//REMOVE i from  dataResponse.Date = extractDate(currentTime.AddDate(0, 0, i))
//TO CALCULTE DATA EVERY DAY instead of every minute

//InsertLatestRatio will insert the latest currency ratios into database
func (con CurrConnection) InsertLatestRatio() error {
	latestPath := buildPath(apiLatestEndpoint)
	bodyBites, err := runClient(latestPath)
	if err != nil {
		return err
	}
	var dataResponse RatesResponse
	if err := json.Unmarshal(bodyBites, &dataResponse); err != nil {
		return err
	}

	currentTime := time.Now()

	dataResponse.Date = extractDate(currentTime.AddDate(0, 0, i))
	i = i + 1
	if _, err := con.InsertRatio(dataResponse); err != nil {
		return err
	}
	return nil
}

//REMOVE +i in	startDate := extractDate(currentTime.AddDate(0, 0, -longTermMedian+1+i)) and 	endDate := extractDate(currentTime.AddDate(0, 0, i))
//TO GET DATA EVERY DAY insted every minute and calculate it like its every day

//MedianControler return long term median for longTermMedian days, and short median for shortTermMedian days
func (con CurrConnection) MedianControler(baseCurrency string, tradeCurrency string, longTermMedian int, shortTermMedian int) (float32, float32, error) {
	if longTermMedian < shortTermMedian {
		return 0, 0, errors.New("error: long term median must be larger then short term median")
	}
	if longTermMedian > initializingRatesForDays {
		return 0, 0, fmt.Errorf("error: long term median can not be larger then %d", initializingRatesForDays)
	}
	symbols, err := con.GetAllCurrencySymbols()

	if err != nil {
		return 0, 0, err
	}
	if longTermMedian < shortTermMedian {
		return 0, 0, fmt.Errorf("error: longTermMedian must be a larger number then shortTermMedian")
	}
	baseID, ok := symbols[baseCurrency]
	if ok == false {
		return 0, 0, fmt.Errorf("error: there is no currency  symbol '%s' in database", baseCurrency)
	}
	tradeID, ok := symbols[tradeCurrency]
	if ok == false {

		return 0, 0, fmt.Errorf("error: there is no currency  symbol '%s' in database", tradeCurrency)
	}
	currentTime := time.Now()
	startDate := extractDate(currentTime.AddDate(0, 0, -longTermMedian+1+i))
	endDate := extractDate(currentTime.AddDate(0, 0, i))
	querryString := fmt.Sprintf("WITH normalized_first as (select rate from currencySQL.rates where id1=%d and id2=%d  and creation_date BETWEEN '%s' AND '%s'  union select 1/rate from currencySQL.rates where id1=%d and id2=%d and creation_date BETWEEN '%s' AND '%s') select avg(rate) from normalized_first;", baseID, tradeID, startDate, endDate, tradeID, baseID, startDate, endDate)

	result, err := con.DBConnect.Query(querryString)
	if err != nil {
		return 0, 0, err
	}
	defer result.Close()
	var longMedian float32
	for result.Next() {
		err = result.Scan(&longMedian)
		if err != nil {
			return 0, 0, err
		}
	}

	startDate = extractDate(currentTime.AddDate(0, 0, -shortTermMedian+1+1))
	querryString = fmt.Sprintf("WITH normalized_first as (select rate from currencySQL.rates where id1=%d and id2=%d  and creation_date BETWEEN '%s' AND '%s'  union select 1/rate from currencySQL.rates where id2=%d and id1=%d and creation_date BETWEEN '%s' AND '%s') select avg(rate) from normalized_first;", baseID, tradeID, startDate, endDate, baseID, tradeID, startDate, endDate)

	result1, err := con.DBConnect.Query(querryString)
	if err != nil {
		return 0, 0, err
	}
	defer result1.Close()
	var shortMedian float32
	for result1.Next() {
		err = result1.Scan(&shortMedian)
		if err != nil {
			return 0, 0, err
		}
	}

	return longMedian, shortMedian, nil
}

//convertResponseRatesToNewBase takes two args, RatesResponse struct and newBase string
//it will set RatesResponse.Base to newBase and will convert Rates so they will
//be corresponding trade rates for new base
//This function will not put old base into new Rates
func convertResponseRatesToNewBase(rates RatesResponse, newBase string) (RatesResponse, error) {
	if _, ok := rates.Rates[newBase]; ok == false {
		return rates, fmt.Errorf("error occured; newBase is not a valid currency symbol")
	}
	oldRate := rates.Rates[newBase]
	rates.BaseCurrency = newBase
	for k, v := range rates.Rates {
		rates.Rates[k] = v / oldRate
	}

	return rates, nil
}

func (con CurrConnection) UpdateDaily() {
	wg := &sync.WaitGroup{}
	var updateFuncError error
	updateFunc := func() {
		updateFuncError = con.InsertLatestRatio()
		if updateFuncError != nil {
			wg.Done()
		} else {
			initializingRatesForDays++
		}
		fmt.Printf("Database updated\n")
	}
startUpdate:
	wg.Add(1)
	c := cron.New()
	c.AddFunc("* * * * *", updateFunc)
	c.Start()
	wg.Wait()

	c.Stop()
	fmt.Printf("Error occured: Unable to update databse daily\n%s\n", updateFuncError.Error())
	updateFuncError = nil
	goto startUpdate
}

func extractDate(date time.Time) string {
	return strings.Split(date.Format(dateFormat), " ")[0]
}

func buildPath(endpoint string) string {
	return fmt.Sprintf("%s/%s%s", apiBasePath, endpoint, accessKey)
}
