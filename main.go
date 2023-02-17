package main

import (
	"auth_service/bootstrap"
	"log"
)

func main() {

	err := bootstrap.Run()
	if err != nil {
		log.Fatal(err)
	}
}
