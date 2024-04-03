package main

import (
	"log"
	"os"

	"github.com/merliot/temp"
)

//go:generate go run main.go
func main() {
	temp := temp.New("proto", "temp", "proto").(*temp.Temp)
	if err := temp.GenerateUf2s("../.."); err != nil {
		log.Println("Error generating UF2s:", err)
		os.Exit(1)
	}
}
