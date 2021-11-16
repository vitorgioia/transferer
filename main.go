package main

import (
	"log"
	"net/http"
)

const port = ":8000"

func main() {
	store := NewInMemoryTransfererStore()
	server := &TransfererServer{store}

	log.Fatal(http.ListenAndServe(port, server))
}
