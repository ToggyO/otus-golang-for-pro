package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func listen(client TelnetClient) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGTERM)
	signal.Ignore(syscall.SIGINT)
	defer signal.Stop(signalChan)

	quitChan := NewBufferedFinalizationChannel(1)

	go func(tnc TelnetClient) {
		_ = tnc.Send()
		signalChan <- syscall.SIGQUIT
	}(client)

	go func(done *FinalizationChannel, tnc TelnetClient) {
		_ = tnc.Receive()
		done.SafeClose()
	}(quitChan, client)

	select {
	case <-quitChan.C:
		log.Println("...Connection was closed by peer")
	case <-signalChan:
		log.Println("...Connection was closed by client")
	}
}
