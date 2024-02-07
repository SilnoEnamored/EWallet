package model

import "time"

type Transaction struct {
	ID           int64     `json:"id"`
	FromWalletID int64     `json:"from_wallet_id"`
	ToWalletID   int64     `json:"to_wallet_id"`
	Amount       float64   `json:"amount"`
	Time         time.Time `json:"time"`
}
