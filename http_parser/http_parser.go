package http_parser

import (
	"errors"
	"fmt"
	"strings"
)

const CLRF = "\r\n"
const HEADER_BODY_SPLIT = "\r\n\r\n"

type HttpRequest struct {
	Method          string
	Endpoint        string
	QueryParameters map[string]string
	Headers         map[string]string
	Body            string
}

type HttpResponse struct {
	Version       string
	StatusCode    string
	StatusMessage string
	Headers       map[string]string
	Body          string
}

func (hr HttpResponse) Build() []byte {
	header_clrf := []byte(CLRF)
	response := []byte("HTTP")
	// for the status line
	// add /verion
	// add space
	// add status code
	// add status message
	// add CLRF
	response = append(response, []byte(fmt.Sprintf("/%s", hr.Version))...)
	response = append(response, 32)
	response = append(response, []byte(hr.StatusCode)...)
	response = append(response, 32)
	response = append(response, []byte(hr.StatusMessage)...)

	response = append(response, []byte(CLRF)...)

	// for headers
	// {header name}: {header value}
	// CLRF
	// for all headers
	var headers []byte
	for header_name, header_value := range hr.Headers {
		local := []byte(fmt.Sprintf("%s: %s", header_name, header_value))
		headers = append(headers, local...)
		headers = append(headers, header_clrf...)
	}
	response = append(response, headers...)
	response = append(response, []byte(CLRF)...)

	// CLRF
	// body
	response = append(response, []byte(hr.Body)...)
	return response
}

func InternalServiceResponse() HttpResponse {
	return HttpResponse{Version: "1.1", StatusCode: "500", StatusMessage: "Internal Service Error"}
}

func NotFoundResponse() HttpResponse {
	return HttpResponse{Version: "1.1", StatusCode: "404", StatusMessage: "Not Found"}
}

func BadRequest() HttpResponse {
	return HttpResponse{Version: "1.1", StatusCode: "400", StatusMessage: "Bad Request"}
}

func parseHeaders(h string) map[string]string {
	headerMap := make(map[string]string)
	headers := strings.Split(h, CLRF)
	for _, header := range headers {
		header_split := strings.Split(header, ": ")
		if len(header_split) < 2 {
			continue
		}
		header_name := strings.ToLower(header_split[0])
		header_val := strings.TrimSpace(header_split[1])
		headerMap[header_name] = header_val
	}
	return headerMap
}

func parseQueryParameters(endpoint string) (map[string]string, error) {
	// TODO: add an error possibility here for invalid url
	params := make(map[string]string)
	epSplit := strings.Split(endpoint, "?")
	if len(epSplit) == 1 {
		return params, nil
	}
	p := strings.Split(epSplit[1], "&")
	len_params := len(p)
	if len_params%2 != 0 {
		return params, errors.New("Invalid parameters")
	}
	for i := 0; i < len_params; i += 2 {
		params[p[i]] = p[i+1]
	}

	return params, nil
}

func ParseRequest(b []byte) (HttpRequest, error) {
	request := string(b)
	req_split := strings.Split(request, CLRF)
	len_status_line := len(req_split[0])
	status_line := strings.Split(req_split[0], " ")
	method := status_line[0]
	endpoint := status_line[1]
	query_parameters, err := parseQueryParameters(endpoint)
	if err != nil {
		return HttpRequest{}, errors.New("Invalid http request")
	}
	header_and_body := strings.Split(request[len_status_line:], HEADER_BODY_SPLIT)
	header_line := header_and_body[0]
	headers := parseHeaders(header_line)
	var body string
	if len(header_and_body) == 2 {
		body = header_and_body[1]
	} else {
		body = ""
	}
	return HttpRequest{
		Method: string(method), Endpoint: string(endpoint), QueryParameters: query_parameters, Headers: headers, Body: body,
	}, nil

}
