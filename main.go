package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	s := &http.Server{ //TODO: TLS
		Addr:           ":8080", //TODO: use config
		Handler:        newApp(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, //1MB
	}

	log.Printf("Listening... %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleInterrupts(ch chan os.Signal) {
	signal.Notify(ch, os.Interrupt)
	go func() {
		for sig := range ch {
			fmt.Printf("Exiting... %v\n", sig)
			ch = nil
			os.Exit(1)
		}
	}()
}
