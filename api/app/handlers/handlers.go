package api

import (
	"EWallet/api/app/model"
	"EWallet/api/pkg/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Обработчики HTTP-запросов
func CreateWalletHandler(walletStore *store.WalletStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Создаем кошелек с балансом 100
		wallet := &model.Wallet{
			Balance: 100.0,
		}

		// Добавляем кошелек в бд
		err := walletStore.CreateWallet(wallet)
		if err != nil {
			log.Printf("Error creating wallet: %v", err)
			http.Error(w, "Error creating wallet", http.StatusInternalServerError)
			return
		}

		// заголовок + отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wallet)
	}
}

type TransactionRequest struct {
	To     int64   `json:"to"`
	Amount float64 `json:"amount"`
}

func MakeTransactionHandler(walletStore *store.WalletStore, transactionStore *store.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		walletId, err := strconv.ParseInt(vars["walletId"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		var req TransactionRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		senderWallet, err := walletStore.GetWalletByID(walletId)
		if err != nil {
			http.Error(w, "Sender wallet not found", http.StatusNotFound)
			return
		}
		if senderWallet == nil {
			http.Error(w, "Sender wallet not exists", http.StatusNotFound)
			return
		}

		if senderWallet.Balance < req.Amount {
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
			return
		}

		recipientWallet, err := walletStore.GetWalletByID(req.To)
		if err != nil || recipientWallet == nil {
			http.Error(w, "Recipient wallet not found", http.StatusBadRequest)
			return
		}

		senderWallet.Balance -= req.Amount
		recipientWallet.Balance += req.Amount

		transaction := &model.Transaction{
			FromWalletID: walletId,
			ToWalletID:   req.To,
			Amount:       req.Amount,
			Time:         time.Now(),
		}

		err = transactionStore.CreateTransaction(transaction)
		if err != nil {
			log.Printf("Error recording transaction: %v", err)
			http.Error(w, "Error recording transaction", http.StatusInternalServerError)
			return
		}

		err = walletStore.UpdateWallet(senderWallet)
		if err != nil {
			log.Printf("Error updating sender wallet: %v", err)
			http.Error(w, "Error updating sender wallet", http.StatusInternalServerError)
			return
		}

		err = walletStore.UpdateWallet(recipientWallet)
		if err != nil {
			log.Printf("Error updating recipient wallet: %v", err)
			http.Error(w, "Error updating recipient wallet", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Transaction completed successfully"))
	}
}

func GetWalletHandler(walletStore *store.WalletStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		walletId, err := strconv.ParseInt(vars["walletId"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		wallet, err := walletStore.GetWalletByID(walletId)
		if err != nil || wallet == nil {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wallet)
	}
}

func GetTransactionsHandler(transactionStore *store.TransactionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		walletId, err := strconv.ParseInt(vars["walletId"], 10, 64)
		if err != nil {
			http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
			return
		}

		transactions, err := transactionStore.GetTransactionsByWalletID(walletId)
		if err != nil {
			http.Error(w, "Error retrieving transactions", http.StatusInternalServerError)
			return
		}

		if transactions == nil || len(transactions) == 0 {
			http.Error(w, "No transactions found for the wallet", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transactions)
	}
}

// маршруты для обработчиков
func SetupRoutes(walletStore *store.WalletStore, transactionStore *store.TransactionStore) {
	http.HandleFunc("/wallet/create", CreateWalletHandler(walletStore))
	http.HandleFunc("/transaction/make", MakeTransactionHandler(walletStore, transactionStore))
	http.HandleFunc("/wallet/{walletId}", GetWalletHandler(walletStore))
	http.HandleFunc("/transactions/", GetTransactionsHandler(transactionStore))
}
