package main

import(
	"fmt"
	"log"
	"bytes"
	"io"
	"net"
	)


func getLineChannel(f io.ReadCloser) chan string{
	out := make(chan string, 2)

	go func(){
		defer f.Close()
		defer close(out)
		line := ""
	for {
	 data := make([]byte, 8)
	 n,err :=  f.Read(data)


	 if err != nil{
		break
	 }

	data = data[:n]
	i :=bytes.IndexByte(data,'\n')
	if i!=-1{
		line += string(data[:i])
		data = data[i+1:]
		out <- line
		line =""

	}
	line += string(data)
	}
	if(len(line) !=0){
		out <- line
	}
	}()
	return out
}
func main(){
	listenr,err := net.Listen("tcp", ":42069")
	if err != nil{
		log.Fatal("error",err)
	}
	conn, err := listenr.Accept()
	if (err != nil){
		log.Fatal("error : ", err)
	}
	lines := getLineChannel(conn);

	for line :=  range(lines){
		fmt.Println(line)
	}
	
}
