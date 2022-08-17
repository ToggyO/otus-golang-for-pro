package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

var timeout time.Duration

func main() {
	defaultTimeout, _ := time.ParseDuration("10s")
	flag.DurationVar(&timeout, "timeout", defaultTimeout, "Connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("required arguments \"host\" and \"port\" not define")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	address := net.JoinHostPort(host, port)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	log.Printf("...Connected to %s", address)
	defer client.Close()

	listen(client)
}
