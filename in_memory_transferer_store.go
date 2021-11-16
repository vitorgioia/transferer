package main

func NewInMemoryTransfererStore() *InMemoryTransfererStore {
	return &InMemoryTransfererStore{}
}

type InMemoryTransfererStore struct {
	accountList []Account
}

func (i *InMemoryTransfererStore) GetAccounts() []Account {
	return i.accountList
}

func (i *InMemoryTransfererStore) PostAccount(a Account) {
	i.accountList = append(i.accountList, a)
}

func (s *InMemoryTransfererStore) GetAccountBalance(id string) string {
	for _, account := range s.GetAccounts() {
		if account.Id == id {
			return account.Balance
		}
	}

	return ""
}
