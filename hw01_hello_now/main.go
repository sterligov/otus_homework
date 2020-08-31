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

	fmt.Printf("current time: %s\nexact time: %s\n", currentTime, exactTime)
}
