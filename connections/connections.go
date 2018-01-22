package connections

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"net"

	"github.com/nseveryns/webserver/configuration"
	"github.com/nseveryns/webserver/messages"
	"github.com/nseveryns/webserver/provider/page"
	"github.com/nseveryns/webserver/provider/static"
)

type Provider interface {
	ProvidePage(s string) *page.Page
}

type wrapper struct {
	conn     net.Conn
	config   configuration.Configuration
	request  *messages.RequestMessage
	response *messages.ResponseMessage
	provider Provider
}

func Create(conn net.Conn, config configuration.Configuration) *wrapper {
	provider := &static.StaticProvider{config, make(map[string]static.StaticPage)}
	return &wrapper{conn: conn, config: config, provider: provider}
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

func (wrapper *wrapper) HandleConnection() {
	conn := wrapper.conn
	log.Print("Handling new connection from", conn.LocalAddr())
	scanner := bufio.NewScanner(conn)
	scanner.Split(splitMessage)
	for scanner.Scan() {
		wrapper.request = messages.ReadRequestMessage(scanner.Text())
		err := wrapper.processRequest()
		if err != nil {
			log.Fatal("Unable to generate response message", err)
		}
		conn.Write(messages.GenerateResponse(*wrapper.response))
		conn.Close()
		return
	}
}

func (wrapper *wrapper) processRequest() error {
	request := wrapper.request
	response := &messages.ResponseMessage{} //TODO Look into provider creating response
	err := errors.New("Unable to get identify method type.")
	if request.Method == messages.GET {
		response.Version = "HTTP/" + request.Version
		if request.Path == "/" { //If there is no page request
			request.Path = wrapper.config.LandingPage
		}
		page := wrapper.provider.ProvidePage(request.Path)
		response.Code = messages.OK
		response.Message = string(page.Content)
		wrapper.response = response
		err = nil
	}
	return err
}

func (wrapper *wrapper) WrapUp() {
	wrapper.conn.Close()
}
