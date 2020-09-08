package service

import (
	"currency/client"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Handlers interface {
	SetupRoutes()
}

func (con Connection) SetupRoutes() {
	http.HandleFunc("/api/trade/", con.handleAverage)

	http.HandleFunc("/api/inform/", con.handleMovingAverage)

	http.HandleFunc("/api/register/", con.handleRegistration)
	checkError(http.ListenAndServe(":8080", nil))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (con Connection) handleMovingAverage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if ok, _ := regexp.MatchString(RegexApiInform, r.URL.Path); ok == false {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Sprintf("error: wrong path\npath schema: %s\npath example: localhost:8080/api/inform/aleksa/USD/GBP/7/3", RegexApiInform)
		w.Write([]byte(err))
		fmt.Printf("%s\n", err)
		return
	}
	urlPathSegments := strings.Split(r.URL.Path, "/")
	nick := urlPathSegments[3]
	baseCurr := urlPathSegments[4]
	tradeCurr := urlPathSegments[5]
	longMedian, _ := strconv.Atoi(urlPathSegments[6])
	shortMedian, _ := strconv.Atoi(urlPathSegments[7])

	clientConnection := client.ClConnection(con) //{database.Connection{DBConnect: con.DBConnect}}
	client, err := clientConnection.GetClientByName(nick)
	if err != nil {
		errSplit := strings.Split(err.Error(), ":")
		if errSplit[0] == "error" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Printf("%s\n", err.Error())
		return
	}

	currentMedianJSON, err := con.MedianDifferenceReturn(baseCurr, tradeCurr, longMedian, shortMedian)
	if err != nil {
		errSplit := strings.Split(err.Error(), ":")
		if errSplit[0] == "error" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			fmt.Printf("%s\n", err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("%s\n", err.Error())
		return
	}
	var currentMedian MedianAnswer
	if err := json.Unmarshal(currentMedianJSON, &currentMedian); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("%s\n", err.Error())
		return
	}

	go con.MovingAverage(client.Email, currentMedian, baseCurr, tradeCurr, longMedian, shortMedian)
	w.Write([]byte("you will be informed via email.. when i make it work"))
}

func (con Connection) handleRegistration(w http.ResponseWriter, r *http.Request) { // func(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonAnswer, err := con.RegisterNewClientURL(r.URL.Path, RegexApiRegister)
	if err != nil {
		errSegments := strings.Split(err.Error(), ":")
		if errSegments[0] == "error" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("%s\n", err.Error())
	}
	w.Write(jsonAnswer)

}

func (con Connection) handleAverage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonAnswer, err := con.MedianDifferenceReturnURL(r.URL.Path, RegexApiTrade, 3)
	if err != nil {
		errSegments := strings.Split(err.Error(), ":")
		if errSegments[0] == "error" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Printf("%s\n", err.Error())
		return
	}
	w.Write(jsonAnswer)
}
