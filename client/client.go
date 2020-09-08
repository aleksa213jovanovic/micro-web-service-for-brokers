package client

import "currency/database"

type Client struct {
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
}

type ClConnection struct {
	database.Connection
}
