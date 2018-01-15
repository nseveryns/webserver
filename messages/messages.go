package messages

import (
	"fmt"
	"strings"
)

type Method string
type StatusCode int16

const (
	GET  Method = "GET"
	HEAD Method = "HEAD"
	POST Method = "POST"
	NIL  Method = "UNKNOWN"

	OK        StatusCode = 200
	NOT_FOUND StatusCode = 404
)

type RequestMessage struct {
	Method  Method
	Path    string
	Version string
	Header  map[string]string
	Message string
}

type ResponseMessage struct {
	Version string
	Code    StatusCode
	Header  map[string]string
	Message string
}

func ReadRequestMessage(raw string) *RequestMessage {
	fmt.Println(raw)
	split := strings.Split(raw, " ")
	fmt.Println(split[0])
	method, err := getMethod(split[0])
	if err != nil {
		fmt.Println("Invalid Method: ", split[0], err)
	}
	path := split[1]
	fmt.Println(path)
	version := split[2][5:8]
	fmt.Println(version)
	flags := getFlags(raw[strings.Index(raw, "\r\n")+1:])
	message := new(RequestMessage)
	message.Method = method
	message.Path = path
	message.Version = version
	message.Header = flags
	return message
}

func GenerateResponse(response ResponseMessage) []byte {
	sprint := fmt.Sprintf("%s %d %s\n\n%s", response.Version, response.Code, "OK", response.Message)
	fmt.Println(response.Code)
	fmt.Println(sprint)
	return []byte(sprint)
}
