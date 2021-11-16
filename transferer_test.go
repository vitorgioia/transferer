package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubTransfererStore struct {
	accountList []Account
}

func (s *StubTransfererStore) GetAccounts() []Account {
	return s.accountList
}

func (s *StubTransfererStore) PostAccount(a Account) {}

func TestGETAccounts(t *testing.T) {

	t.Run("it returns the accounts as JSON", func(t *testing.T) {
		wantedAccounts := []Account{
			{"xyz", "John", "10.00"},
			{"abc", "Mary", "20.00"},
		}
		store := new(StubTransfererStore)
		store.accountList = wantedAccounts

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server := TransfererServer{store}

		server.ServeHTTP(response, request)

		var got []Account
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response %q into slice of Accounts, '%v'", response.Body, err)
		}
		assertStatus(t, response.Code, http.StatusOK)
		asserContentType(t, response, jsonContentType)

		if !reflect.DeepEqual(got, wantedAccounts) {
			t.Errorf("got %v want %v", got, wantedAccounts)
		}
	})
}

func TestPOSTAccounts(t *testing.T) {

	t.Run("it returns a 201 response", func(t *testing.T) {
		store := new(StubTransfererStore)
		server := TransfererServer{store}

		jsonStr := []byte(`{"id": "xyz", "name": "John", "balance": "0.00"}`)
		request, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonStr))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		wantedAccount := Account{"xyz", "John", "0.00"}

		var got Account
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response %q into slice of Accounts, '%v'", response.Body, err)
		}
		assertStatus(t, response.Code, http.StatusCreated)
		asserContentType(t, response, jsonContentType)

		if got != wantedAccount {
			t.Errorf("got %q want %q", got, wantedAccount)
		}
	})
}

func assertStatus(t testing.TB, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}
}

func asserContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
	}
}
