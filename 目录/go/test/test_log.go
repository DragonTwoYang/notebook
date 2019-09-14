package main

import (
	"os"
	"log"
	"fmt"
)


func main() {
	file, err := os.Create("log.txt")
	if err != nil {

		os.Exit(1)
	}
}