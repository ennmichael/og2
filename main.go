package main

import (
	"fmt"
	"og2/game"
	"og2/rest"
)

func main() {
	store, err := game.InitStore()
	fmt.Println("Init done")
	if err != nil {
		panic(err)
	}
	err = rest.RunServer(store)
	if err != nil {
		panic(err)
	}
}
