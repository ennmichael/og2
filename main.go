package main

import (
	"og2/rest"
)

func main() {
	err := rest.RunServer()
	if err != nil {
		panic(err)
	}
}
