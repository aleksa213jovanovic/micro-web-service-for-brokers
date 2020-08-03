package client

import (
	"log"
	"os"
)

type ClientConnection interface {
	SetupClient()
}

func (con ClConnection) SetupClient() {
	checkError(createClientTable(con))
}

func createClientTable(con ClConnection) error {
	querry := `CREATE TABLE currencySQL.clients(
		id int not null auto_increment,
		nick_name varchar(255) not null,
		email varchar(255) not null,
		PRIMARY KEY(id)
		);`
	_, err := con.DBConnect.Exec(querry)
	return err
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		
	}
}
