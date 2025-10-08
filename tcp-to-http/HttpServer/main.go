package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tcp-to-http/internal/server"
)

const port = 42069

func main(){
	s , err := server.Serve(port)
	if err != nil{
		log.Fatalf("some error in server : %v", err)
	}

	defer s.Close()


	log.Println("server up and running on prt :", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) 
	<- sigChan

	log.Println("server gracefully stopped ")
}