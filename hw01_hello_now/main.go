package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const ntpServer = "0.beevik-ntp.pool.ntp.org"

func main() {
	exactTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.Fatalln(err)
	}

	currentTime := time.Now()

	fmt.Printf("current time: %s\n", currentTime.Round(0))
	fmt.Printf("exact time: %s\n", exactTime.Round(0))
}
