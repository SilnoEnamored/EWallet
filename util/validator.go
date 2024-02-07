package util

import (
	"EWallet/api/app/model"
	"errors"
)

// проверка корректности данных кошелька.
func ValidateWallet(wallet *model.Wallet) error {
	if wallet.Balance < 0 {
		return errors.New("balance cannot be negative")
	}
	return nil
}

// проверка корректности данных транзакции.
func ValidateTransaction(transaction *model.Transaction) error {
	if transaction.Amount <= 0 {
		return errors.New("transaction amount must be greater than 0")
	}
	if transaction.FromWalletID == transaction.ToWalletID {
		return errors.New("sender and recipient wallets must be different")
	}
	return nil
}
