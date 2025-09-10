package main

import(
	"fmt"
	"log"
	"os"
	"bytes"
	"io"
	"time"
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
		time.Sleep(5 * time.Second)
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
	var f,err = os.Open("message.txt")
	if err != nil{
		log.Fatal("error",err)
	}
	
	lines := getLineChannel(f);

	for line :=  range(lines){
		fmt.Println(line)
	}
	
}
