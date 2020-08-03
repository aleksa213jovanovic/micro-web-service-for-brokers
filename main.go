package main

import (
	"currency/client"
	"currency/currency"
	"currency/database"
	"currency/service"
	"log"
)

const (
	driverName     = "mysql"
	dataSourceName = "aleksa:qweQWE1_@tcp(127.0.0.1:3306)/currencySQL"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	con, err := database.NewConnection(driverName, dataSourceName)
	checkError(err)
	currencyCon := currency.CurrConnection{con} 
	currencyCon.SetupCurrency()
	clientCon := client.ClConnection{con} 
	clientCon.SetupClient()
	serverConnection := service.Connection{con}
	serverConnection.SetupRoutes( ) 
}
