package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/nseveryns/webserver/messages"
)

func main() {
	fmt.Print("Starting web server.")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Unable to bind to port.", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept new connection.", err)
		}
		go handleConnection(conn)
		defer conn.Close()
	}

}

func splitMessage(data []byte, atEOF bool) (int, []byte, error) {
	length := len(data)
	if atEOF && length == 0 {
		return 0, data, nil
	}
	if length < 4 {
		return 0, nil, nil
	}
	cropLength := length - 4
	d := data[cropLength:]
	if bytes.Equal(d[:], []byte("\r\n\r\n")) {
		return len(data), data[:cropLength], nil
	}
	return 0, nil, nil
}

func handleConnection(conn net.Conn) {
	fmt.Println("Handling new connection from", conn.LocalAddr())
	scanner := bufio.NewScanner(conn)
	scanner.Split(splitMessage)
	for scanner.Scan() {
		message := messages.ReadRequestMessage(scanner.Text())
		response, err := processRequest(*message)
		if err != nil {
			fmt.Println("Unable to generate response message", err)
		}
		conn.Write(messages.GenerateResponse(*response))
		conn.Close()
		return
	}
}

func processRequest(request messages.RequestMessage) (*messages.ResponseMessage, error) {
	response := new(messages.ResponseMessage)
	err := errors.New("Unable to get identify method type.")
	if request.Method == messages.GET {
		response.Version = "HTTP/" + request.Version
		data, ioerror := ioutil.ReadFile("." + request.Path)
		if ioerror != nil {
			response.Code = messages.NOT_FOUND
			data = []byte("Not found!")
			fmt.
				Println(request.Path)
		}
		response.Code = messages.OK
		response.Message = string(data)
		err = nil
	}
	return response, err
}
