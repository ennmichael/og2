package main

import (
	"og2/rest"
	"og2/store"
)

func main() {
	store, err := store.Init()
	if err != nil {
		panic(err)
	}
	err = rest.RunServer(store)
	if err != nil {
		panic(err)
	}
}
