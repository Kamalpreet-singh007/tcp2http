package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tcp-to-http/internal/server"
	"tcp-to-http/internal/requests"
	"tcp-to-http/internal/response"
	"fmt"
)

const port = 42069

func resposnd400 ()[]byte{
	return []byte (`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}
func resposnd500 ()[]byte{
	return []byte (`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}
func resposnd200 ()[]byte{
	return []byte (`  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}
func main(){
	s , err := server.Serve(port, func(w *response.Writer, req *request.Request){

		h := response.GetDefaulHeaders(0)	
		body := resposnd200()
		status :=response.StatusOK

		if req.RequestLine.RequestTarget == "/yourproblem"{	
			body = resposnd400()
			status =response.StatusBad
		}else if  req.RequestLine.RequestTarget == "/myproblem"{
			body = resposnd500()
			status =response.StatusBad 
		}

		w.WriteStatusLine(status)
		h.Replace("content-length", fmt.Sprintf("%d",len(body)))
		h.Replace("content-type","text/html")
		w.WriteHeaders(*h)
		w.WriteBody(body)
	})
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