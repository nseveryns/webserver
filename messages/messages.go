package messages

import (
	"fmt"
	"log"
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

//TODO Less hard coding would be nice
func ReadRequestMessage(raw string) *RequestMessage {
	split := strings.Split(raw, " ")
	method, err := getMethod(split[0])
	if err != nil {
		log.Println("Invalid Method: ", split[0], err)
	}
	path := split[1]
	version := split[2][5:8]
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
	return []byte(sprint)
}
