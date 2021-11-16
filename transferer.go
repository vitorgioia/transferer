package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const jsonContentType = "application/json"

type Account struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Balance string `json:"balance"`
}

type TransfererStore interface {
	GetAccounts() []Account
	PostAccount(a Account)
}

type TransfererServer struct {
	store TransfererStore
}

func (t *TransfererServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	switch r.Method {
	case http.MethodGet:
		t.fetchAccountList(w)
	case http.MethodPost:
		t.createAccount(w, r)
	}
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
