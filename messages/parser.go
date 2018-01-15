package messages

import (
	"errors"
	"fmt"
	"strings"
)

func getMethod(unparsed string) (Method, error) {
	var method = NIL
	switch unparsed {
	case "GET":
		method = GET
		break
	case "HEAD":
		method = HEAD
		break
	case "POST":
		method = POST
	}
	var err error
	if method == NIL {
		err = errors.New("Method was not valid")
		fmt.Println(unparsed)
	}
	return method, err
}

func getFlags(unparsed string) map[string]string {
	flags := make(map[string]string)
	lines := strings.Split(unparsed, "\r\n")
	for _, element := range lines {
		pair := strings.Split(element, ": ")
		if len(pair) > 2 {
			fmt.Println("Invalid content flag.", element)
			break
		}
		flags[pair[0]] = pair[1]
	}
	return flags
}
