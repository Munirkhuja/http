package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Munirkhuja/http/pkg/server"
)
const crlf="\r\n"
func main() {
	host := "0.0.0.0"
	port := "9999"
	if err:=execute(host,port);err!=nil {
		os.Exit(1)
	}
}
func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host,port))
	srv.Register("/",func(conn net.Conn) {
		body:="Welcome to our web-site"
		_,err=conn.Write([]byte(httpDefaultHeader(body)))
		if err!=nil {
			log.Print(err)
		}
	})
	srv.Register("/about",func(conn net.Conn) {
		body:="About Golang Academy"		
		_,err=conn.Write([]byte(httpDefaultHeader(body)))
		if err!=nil {
			log.Print(err)
		}
	})
	return srv.Start()
}
func httpDefaultHeader(body string)string{
	result:="HTTP/1.1 200 OK"+crlf+
	"Content-Length:"+strconv.Itoa(len(body))+crlf+
	"Content-Type: text/html"+crlf+
	"Connection: close"+crlf+
	crlf+
	body
	return result
}