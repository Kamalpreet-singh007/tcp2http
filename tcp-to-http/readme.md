# ⚙️ TCP to HTTP Server

This is a simple **TCP-to-HTTP server** built in Go. It listens for raw TCP connections and responds with custom **HTTP-style responses** depending on the requested path. It’s a small, educational project to understand how HTTP can be simulated manually over TCP sockets.

To run the project, go to the root directory and execute:

```bash
go run HttpServer/main.go


# Returns 400 Bad Request
curl -v http://localhost:42069/yourproblem

# Returns 500 Internal Server Error
curl -v http://localhost:42069/myproblem

# Returns 200 OK
curl -v http://localhost:42069/anythingelse

# Test chunked transfer encoding
echo -e "GET /httpbin/stream/100 HTTP/Host: localhost\r\nConnection: close\r\n\r\n" | nc localhost 42069
```
