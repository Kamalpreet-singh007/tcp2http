package request

import(
	"io"
	"fmt"
	"bytes"
	"errors"
	"tcp/internal/headers"
)
type parserState string

const (
	StateInit parserState = "init"
	StateDone  parserState = "done"
	StateHeaders  parserState = "headers"
	StateError  parserState = "error"
)

type RequestLine struct{
	HttpVersion string 
	RequestTarget string
	Method string
}
type Request struct{
	RequestLine RequestLine
	Headers *headers.Headers
	state parserState
	
}

// func (r *RequestLine) ValidHTTP() bool{
// 	return r.HttpVersion == "HTTP/1.1"
// }
func newRequest() *Request{
	return &Request{
		state :StateInit,
		Headers: headers.NewHeaders(),
	}
}
var Error_Request_In_Error_State = fmt.Errorf("request in error state")
var ERROR_UNSUPPORTED_HTTP_VERSION =fmt.Errorf("unsupported http-version")
var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var SEPRATOR = []byte("\r\n")

func parseRequestLine(b []byte )(*RequestLine, int , error){
	idx:= bytes.Index(b, SEPRATOR)
	if idx == -1 {
	 	return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx +len(SEPRATOR)
	

	parts := bytes.Split(startLine,[]byte(" "))
	if len(parts)!= 3{
		return nil,0,ERROR_MALFORMED_REQUEST_LINE
	}
	httpParts :=bytes.Split(parts[2],[]byte("/"))
	if len(httpParts)!= 2 || string(httpParts[0]) !="HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}
	rl :=  &RequestLine{
		Method:         string(parts[0]),
		RequestTarget:  string(parts[1]),
		HttpVersion: 	 string(httpParts[1]),
	}

	
	return rl, read, nil
}


func (r* Request) parse(data []byte ) (int, error){
		read:= 0	
		outer:
	for{
		switch r.state {
		case StateError:
			return 0, Error_Request_In_Error_State
		case StateInit:

	
			rl, n, 	err := parseRequestLine(data[read:])
			if err!= nil {
				r.state = StateError
			}
			if n == 0{
				break outer
			}
			 
			r.RequestLine = *rl
			read+=n

			r.state = StateHeaders
		case StateHeaders:
			
			n, done, err := r.Headers.Parse(data[read:])

			if err!= nil{
				return 0,  err
			}
			if n == 0{
				break outer
			}
			read+= n	
			if done {
				r.state = StateDone
			}
			
		case StateDone:
			return 0, nil
		default:
			panic("we did something wrong")
		}
	}
	return read, nil
}
func (r* Request) done() bool{
	return (r.state == StateDone|| r.state==StateError)
}



func RequestFromReader(reader io.Reader)(*Request, error){
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen :=0
	// bufIdx:=0
	for !request.done(){
		n ,err :=reader.Read(buf[bufLen:])
		if err!= nil{
			if errors.Is(err, io.EOF){
				request.state = StateDone
				break
			}
			return nil, err
		}
		bufLen+=n

		readN , err := request.parse(buf[:bufLen])
		if err!= nil{
			return nil , err
		}
		
		copy(buf, buf[readN :bufLen])
		bufLen -= readN

	}
	return request ,nil
} 