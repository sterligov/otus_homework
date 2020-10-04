package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout *time.Duration

func init() {
	timeout = flag.Duration("timeout", 100*time.Second, "")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatalln("Expected 2 arguments: host and port")
	}

	hostport := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(hostport, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	cancel := make(chan error)

	go func() {
		if err := client.Receive(); err != nil {
			cancel <- err
		}
		log.Println("...Connection was closed by peer")
		cancel <- nil
	}()

	go func() {
		if err := client.Send(); err != nil {
			cancel <- err
		}
		log.Println("...EOF")
		cancel <- nil
	}()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		<-sigs
		log.Fatalln("...Shutdown")
		cancel <- nil
	}()

	if err := <-cancel; err != nil {
		log.Fatalln(err)
	}
}
