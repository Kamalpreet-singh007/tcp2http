package main

import(
	"fmt"
	"log"
	// "bytes"
	// "io"
	"net"
	"tcp/internal/requests"
	)


func main(){
	listenr,err := net.Listen("tcp", ":42069")
	if err != nil{
		log.Fatal("error",err)
	}
	for {

		conn, err := listenr.Accept()
		if (err != nil){
			log.Fatal("error : ", err)
		}
		
		
		r,_ := request.RequestFromReader(conn)
		
		fmt.Printf("Request line : \n")
		fmt.Printf("-Method: %s\n",r.RequestLine.Method )
		fmt.Printf("-Target  : %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("-Version : %s\n", r.RequestLine.HttpVersion)

	} 
	
}
