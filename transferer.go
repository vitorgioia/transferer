package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const jsonContentType = "application/json"

type Account struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Balance string `json:"balance"`
}

type AccountBalance struct {
	Balance string `json:"balance"`
}

type TransfererStore interface {
	GetAccounts() []Account
	PostAccount(a Account)
	GetAccountBalance(id string) string
}

type TransfererServer struct {
	store TransfererStore
}

func (t *TransfererServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	router := mux.NewRouter()
	router.HandleFunc("/accounts", t.AccountsHandler).Methods("GET", "POST")
	router.HandleFunc("/accounts/{account_id}/balance", t.BalanceHandler).Methods("GET")
	router.ServeHTTP(w, r)
}

func (t *TransfererServer) AccountsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.fetchAccountList(w)
	case http.MethodPost:
		t.createAccount(w, r)
	}
}

func (t *TransfererServer) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["account_id"]
	amount := t.store.GetAccountBalance(accountId)

	if amount == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(AccountBalance{amount})
}

func (t *TransfererServer) fetchAccountList(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(t.store.GetAccounts())
}

func (t *TransfererServer) createAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	var a Account
	err := json.NewDecoder(r.Body).Decode(&a)

	if err != nil {
		log.Fatalf("Unable to parse request %q into an Account, '%v'", r.Body, err)
	}
	t.store.PostAccount(a)

	json.NewEncoder(w).Encode(a)
}
