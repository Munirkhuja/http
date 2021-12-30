package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type HandleFunc func(conn net.Conn)

type Server struct{
	addr string
	mu sync.RWMutex
	handlers map[string]HandleFunc
}

const ErrRequestHadNotEndLine="request had not end line"
const ErrRequestFormatError="request formet error"

func NewServer(addr string) *Server  {
	return &Server{addr: addr,handlers: make(map[string] HandleFunc)}
}
func (s *Server)Register(path string,handler HandleFunc){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path]=handler
}
func (s *Server)Start() error{
	log.Print(s.addr)
	listener,err:=net.Listen("tcp",s.addr)
	if err!=nil {
		log.Print(err)
		return err	
	}
	defer func ()  {
		if cerr:=listener.Close();cerr!=nil{
			if err==nil {
				err=cerr
				return
			}
			log.Print(cerr)
		}
	}()
	for{
		conn,aerr:=listener.Accept()
		if aerr!=nil {
			log.Print(aerr)
			continue
		}
		go s.handler(conn)
	}
}
func (s *Server)handler(conn net.Conn)(err error){
	defer func() {
		cerr:=conn.Close
		if cerr!=nil{
			if err==nil{
				err=errors.New("connection close")
				return
			}
			log.Print(err)
		}
	}()
	buf:=make([]byte,4096)
	n,err:=conn.Read(buf)
	if err==io.EOF {
		log.Printf("%s",buf[:n])
		return nil
		
	}
	if err!=nil {
		return err
	}
	log.Printf("%s",buf[:n])
	data:=buf[:n]
	requestlineDeLim:=[]byte{'\r','\n'}
	requestLineEnd:=bytes.Index(data,requestlineDeLim)
	if requestLineEnd==-1 {
		err=errors.New(ErrRequestHadNotEndLine)
		return err
	}
	requestLine:=string(data[:requestLineEnd])
	parts:=strings.Split(requestLine," ")
	if len(parts)!=3 {
		err=errors.New(ErrRequestFormatError)
		return err		
	}
	path:=parts[1]
	if path!="" {
		log.Print(path)
	}
	if handlerFunc,found:=s.handlers[path];found {
		handlerFunc(conn)
	}
	return err
}