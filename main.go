package main

import (
	"fmt"
	"go_http/http_parser"
	"net"
	"os"
	"strings"
)

type HttpMethod func(http_parser.HttpRequest) http_parser.HttpResponse
type Server struct {
	tcpListener net.Listener
	routes      map[string]map[string]HttpMethod
}

func (s *Server) add_method(method_type, route string, f HttpMethod) {
	// if the method family is in the map, add it
	m, ok := s.routes[method_type]
	if !ok {
		// otherwise create a new map
		new_m := make(map[string]HttpMethod)
		new_m[route] = f
		s.routes[method_type] = new_m
	}
	m[route] = f
}

func NewServer() Server {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Unable to bind to prot 4221")
		os.Exit(1)
	}
	routes := make(map[string]map[string]HttpMethod)
	return Server{tcpListener: l, routes: routes}
}

func (s *Server) run() {
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		// ship off the conn to a go routine
		// the response is written to the conn in the goroutine
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := c.Read(buffer)
		if err != nil {
			response := http_parser.InternalServiceResponse()
			c.Write(response.Build())
			continue
		}

		buffer = buffer[:n]
		request, err := http_parser.ParseRequest(buffer)
		if err != nil {
			response := http_parser.BadRequest()
			c.Write(response.Build())
		}
		keep_open := true
		conn_header, ok := request.Headers["connection"]
		if conn_header == "close" {
			keep_open = false
		}
		var response http_parser.HttpResponse
		methods, ok := s.routes[request.Method]

		if !ok {
			// if there are no methods for the request type
			// 404 and continue
			response = http_parser.NotFoundResponse()
			c.Write(response.Build())
			continue
		}

		// parse the route
		route := ""
		var http_f HttpMethod
		for m, f := range methods {
			if strings.HasPrefix(request.Endpoint, m) && len(m) > len(route) {
				route = m
				http_f = f
			}
		}
		response = http_f(request)

		c.Write(response.Build())

		// if the connection: closed header is present, close the connection by returning
		if !keep_open {
			c.Close()
			return
		}

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
