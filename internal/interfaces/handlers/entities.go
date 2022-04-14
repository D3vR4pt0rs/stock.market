package handlers

import "github.com/dgrijalva/jwt-go/v4"

type Claims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type UpdateBalanceBody struct {
	Amount float64 `json:"amount"`
}

type StocksBody struct {
	Ticker string `json:"ticker"`
	Amount int    `json:"amount"`
}
