package client

import (
	"fmt"
)

type ClientRepository interface {
	InsertClientRegistration(nickName, email string) error
	GetClientByEmail(email string) (Client, error)
	GetClientByName(nickName string) (Client, error)
	GetClientList() ([]Client, error)
}

func (con ClConnection) InsertClientRegistration(nickName, email string) error {
	_, err := con.DBConnect.Exec(`INSERT INTO currencySQL.clients (nick_name, email) VALUES(?,?);`, nickName, email)
	return err
}
func (con ClConnection) GetClientByEmail(email string) (Client, error) {
	result, err := con.DBConnect.Query(`SELECT nick_name FROM currencySQL.clients WHERE email=?`, email)
	if err != nil {
		return Client{}, err
	}
	defer result.Close()
	var nick string
	for result.Next() {
		result.Scan(&nick)
	}
	if nick == "" {
		return Client{}, fmt.Errorf("error: there is no client with email '%s' in our database", email)
	} else {
		return Client{NickName: nick, Email: email}, nil
	}
}

func (con ClConnection) GetClientByName(nickName string) (Client, error) {
	result, err := con.DBConnect.Query(`SELECT email FROM currencySQL.clients WHERE nick_name=?`, nickName)
	if err != nil {
		return Client{}, err
	}
	defer result.Close()
	var email string
	for result.Next() {
		result.Scan(&email)
	}
	if email == "" {
		return Client{}, fmt.Errorf("error: there is no client with name '%s' in our database", nickName)
	} else {
		return Client{NickName: nickName, Email: email}, nil
	}
}

func (con ClConnection) GetClientList() ([]Client, error) {
	querry := fmt.Sprintf("SELECT nick_name, email FROM currencySQL.clients")
	result, err := con.DBConnect.Query(querry)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	clients := make([]Client, 0)
	var client Client
	for result.Next() {
		result.Scan(&client.NickName, &client.Email)
		clients = append(clients, client)
	}
	return clients, nil
}
