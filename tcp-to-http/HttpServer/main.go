package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"strings"
	"net/http"
	"crypto/sha256"
	
	"tcp-to-http/internal/response"
	"tcp-to-http/internal/server"
	"tcp-to-http/internal/requests"
	"tcp-to-http/internal/headers"
)
const port = 42069

func toStr(bytes []byte) string{
	out := ""

	for _, b := range bytes{
		out+= fmt.Sprintf("%02x",b)
	}
	return out
}

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
		}else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/"){
			target := req.RequestLine.RequestTarget
			res ,err := http.Get("https://httpbin.org/"+target[len("/httpbin/"):])
			if err!= nil{
				body = resposnd500()
				status = response.StatusInternalServerError
			}else{
				w.WriteStatusLine(response.StatusOK)
				h.Delete("content-length")
				h.Set("transfer-encoding", "chunked")
				h.Set("Trailer","X-Content-SHA256")
				h.Set("Trailer","X-Content-Length")
				h.Replace("content-type","text/plain")
				w.WriteHeaders(*h)
				fullBody := []byte{}
				for{
					data := make([]byte,32)
					n,err := res.Body.Read(data)
					if err != nil{
						break
					}
					w.WriteBody([]byte(fmt.Sprintf("%x\r\n",n)))
					w.WriteBody(data[:n])
					fullBody = append(fullBody, data[:n]...)
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n\r\n"))
				trailer := headers.NewHeaders()
				out :=	sha256.Sum256(fullBody)
				trailer.Set("X-Content-SHA256",toStr(out[:]))
				trailer.Set("X-Content-Length", fmt.Sprintf("%d",len(fullBody)) )
				w.WriteHeaders(*trailer)
				w.WriteBody([]byte("\r\n"))
				return
			}
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