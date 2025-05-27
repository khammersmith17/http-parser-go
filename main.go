package main

import (
	"fmt"
	"go_http/http_parser"
	"net"
	"os"
)

type Server struct {
	tcpListener net.Listener
}

func NewServer() Server {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Unable to bind to prot 4221")
		os.Exit(1)
	}
	return Server{tcpListener: l}
}

func (s *Server) run() {
	conn, err := s.tcpListener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// ship off the conn to a go routine
	// the response is written to the conn in the goroutine
	go handleConnection(conn)
}

func handleConnection(c net.Conn) {
	defer c.Close()
	for {
		buffer := make([]byte, 1024)
		n, err := c.Read(buffer)
		if err != nil {
			response := http_parser.HttpResponse{Version: "1.1", StatusCode: "500", StatusMessage: "Internal Service Error"}
			c.Write(response.Build())
			continue
		}

		buffer = buffer[:n]
		request := http_parser.ParseRequest(buffer)
		var response http_parser.HttpResponse

		switch request.Method {
		case "GET":
			response = HandleGetRequest(request)
		case "POST":
			response = HandlePostRequest(request)
		}

		c.Write(response.Build())
	}
}

func HandlePostRequest(r http_parser.HttpRequest) http_parser.HttpResponse {
	var response http_parser.HttpResponse
	return response
}

func HandleGetRequest(r http_parser.HttpRequest) http_parser.HttpResponse {
	var response http_parser.HttpResponse
	return response
}

func main() {
	server := NewServer()
	server.run()
}
