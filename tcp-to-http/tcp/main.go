package main

import(
	"fmt"
	"log"
	// "bytes"
	// "io"
	"net"
	"tcp-to-http/internal/requests"
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
		fmt.Printf("Headers:\n")
		r.Headers.ForEach(func(n,v string){
		fmt.Printf("- %s : %s\n", n,v)
		
	})
	fmt.Printf("Body :\n")
	fmt.Printf("-%s\n", r.Body)
	} 
	
}
