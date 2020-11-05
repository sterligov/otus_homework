package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("expected filename")
	}

	if err := GenerateValidators(os.Args[1]); err != nil {
		log.Fatalln(err)
	}
}
