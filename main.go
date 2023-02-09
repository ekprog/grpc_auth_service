package main

import (
	"Portfolio_Nodes/bootstrap"
	"log"
)

func main() {
	err := bootstrap.Run()
	if err != nil {
		log.Fatal(err)
	}
}
