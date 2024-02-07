package store

import (
	"EWallet/api/app/model"
	"errors"
	"github.com/go-pg/pg/v10"
	"log"
)

type WalletStore struct {
	db *pg.DB
}

func NewWalletStore(db *pg.DB) *WalletStore {
	return &WalletStore{
		db: db,
	}
}

// добавляет новый кошелек в БД
func (store *WalletStore) CreateWallet(wallet *model.Wallet) error {
	_, err := store.db.Model(wallet).Insert()
	return err
}

// возвращает кошелек по его ID
func (store *WalletStore) GetWalletByID(walletID int64) (*model.Wallet, error) {
	wallet := &model.Wallet{ID: walletID}
	err := store.db.Model(wallet).WherePK().Select()
	if err != nil {
		if errors.Is(err.(pg.Error), pg.ErrNoRows) {
			return nil, nil
		}
		log.Print("err", err)
		return nil, err
	}
	return wallet, nil
}

// обновляет баланс кошелька в БД
func (store *WalletStore) UpdateWallet(wallet *model.Wallet) error {
	_, err := store.db.Model(wallet).WherePK().Update()
	return err
}
