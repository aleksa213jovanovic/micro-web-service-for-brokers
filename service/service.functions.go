package service

import (
	"currency/client"
	"currency/currency"
	"currency/database"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/robfig/cron"
)

type ServiceRepository interface {
	RegisterNewClientURL(urlPath, regexCheck string) ([]byte, error)
	RegisterNewClient(nickName, email string) ([]byte, error)
	MedianDifferenceReturnURL(urlPath string, regexCheck string, startIndex int) ([]byte, error)
	MedianDifferenceReturn(baseCurr string, tradeCurr string, longMedian int, shortMedian int) ([]byte, error)
	MovingAverage(email string, currentState MedianAnswer, baseCurr string, tradeCurr string, longMedian int, shortMedian int)
}

func (con Connection) RegisterNewClientURL(urlPath, regexCheck string) ([]byte, error) {

	if ok, _ := regexp.MatchString(regexCheck, urlPath); ok == false {
		return nil, fmt.Errorf("error: wrong path\npath schema: %s\npath example: localhost:8080/api/register/nick/email_123@gmail.com", RegexApiRegister)
	}
	urlPathSegments := strings.Split(urlPath, "/")
	nickName := urlPathSegments[3]
	email := urlPathSegments[4]
	return con.RegisterNewClient(nickName, email)
}

func (con Connection) RegisterNewClient(nickName, email string) ([]byte, error) {
	type Answer struct {
		Status string `json:"status"`
	}
	connection := client.ClConnection{database.Connection{DBConnect: con.DBConnect}}
	if err := connection.InsertClientRegistration(nickName, email); err != nil {
		a := Answer{Status: "Not registerd!"}
		ret, _ := json.Marshal(a)
		return ret, err
	}
	a := Answer{Status: "Registerd!"}
	return json.Marshal(a)
}

func (con Connection) MedianDifferenceReturnURL(urlPath string, regexCheck string, startIndex int) ([]byte, error) {

	if ok, _ := regexp.MatchString(regexCheck, urlPath); ok == false {
		return nil, fmt.Errorf("error: wrong path\npath schema: %s\npath example: localhost:8080/api/trade/USD/RSD/5/2", RegexApiTrade)

	}

	urlPathSegments := strings.Split(urlPath, "/")
	baseCurr := strings.ToUpper(urlPathSegments[startIndex])
	tradeCurr := strings.ToUpper(urlPathSegments[startIndex+1])
	longMedian, err := strconv.Atoi(urlPathSegments[startIndex+2])
	if err != nil {
		return nil, err
	}
	shortMedian, err := strconv.Atoi(urlPathSegments[startIndex+3])
	if err != nil {
		return nil, err
	}
	return con.MedianDifferenceReturn(baseCurr, tradeCurr, longMedian, shortMedian)
}

func (con Connection) MedianDifferenceReturn(baseCurr string, tradeCurr string, longMedian int, shortMedian int) ([]byte, error) {
	baseCurr = strings.ToUpper(baseCurr)
	tradeCurr = strings.ToUpper(tradeCurr)

	connection := currency.CurrConnection{database.Connection{DBConnect: con.DBConnect}}
	l, s, err := connection.MedianControler(baseCurr, tradeCurr, longMedian, shortMedian)
	if err != nil {
		return nil, err
	}

	var a MedianAnswer
	tmp := l - s
	if l < s {
		a = MedianAnswer{Trade: "SELL", MedianDifference: tmp}
	} else {
		a = MedianAnswer{Trade: "BUY", MedianDifference: tmp}
	}
	return json.Marshal(a)
}

func (con Connection) MovingAverage(email string, currentState MedianAnswer, baseCurr string, tradeCurr string, longMedian int, shortMedian int) {

	wg1 := &sync.WaitGroup{}

	var newState MedianAnswer
	calculateNewMedian := func() {
		medianJSON, err := con.MedianDifferenceReturn(baseCurr, tradeCurr, longMedian, shortMedian)
		if err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(medianJSON, &newState); err != nil {
			log.Fatal(err)
		}
		if currentState.Trade != newState.Trade {
			wg1.Done()
		} else {
			fmt.Printf("still same trading state %s, median difference=%f\n", newState.Trade, newState.MedianDifference)
		}
	}
	wg1.Add(1)
	c1 := cron.New()
	c1.AddFunc("@daily", calculateNewMedian)
	c1.Start()

	wg1.Wait()
	c1.Stop()
	fmt.Printf("previous state %s and median difference %f\nnew state %s and median difference %f", currentState.Trade, currentState.MedianDifference, newState.Trade, newState.MedianDifference)
}
