package store

import (
	"EWallet/api/app/model"
	"errors"
	"github.com/go-pg/pg/v10"
)

type TransactionStore struct {
	db *pg.DB
}

func NewTransactionStore(db *pg.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
	}
}

// CreateTransaction добавляет новую транзакцию в БД
func (store *TransactionStore) CreateTransaction(transaction *model.Transaction) error {
	_, err := store.db.Model(transaction).Insert()
	return err
}

// GetTransactionsByWalletID возвращает транзакции по ID кошелька
func (store *TransactionStore) GetTransactionsByWalletID(walletID int64) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := store.db.Model(&transactions).Where("from_wallet_id = ? or to_wallet_id = ?", walletID, walletID).Select()
	if err != nil {
		if errors.Is(err.(pg.Error), pg.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return transactions, err
}
