package server

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	writer          http.ResponseWriter
	request         *http.Request
	destinationHost *string
}

func NewConnection(writer http.ResponseWriter, request *http.Request, destinationHost *string) (*Connection, error) {
	return &Connection{
		writer:          writer,
		request:         request,
		destinationHost: destinationHost,
	}, nil
}

func (c *Connection) Run() error {
	conn, err := upgrader.Upgrade(c.writer, c.request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return err
	}
	defer conn.Close()

	destDialer := websocket.Dialer{}
	destHeaders := http.Header{}

	// Copy selected headers from the original request
	headersToCopy := map[string]bool{"Authorization": true}
	for k, v := range c.request.Header {
		if _, ok := headersToCopy[k]; ok {
			destHeaders.Set(k, v[0])
		}
	}

	queryParams := c.request.URL.RawQuery
	destURL := url.URL{Scheme: "ws", Host: *c.destinationHost, Path: c.request.URL.Path, RawQuery: queryParams}

	var destConn *websocket.Conn
	var destErr error

	// Retry logic
	for i := 0; i < 3; i++ {
		destConn, _, destErr = destDialer.Dial(destURL.String(), destHeaders)
		if destErr == nil {
			break
		}
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if destErr != nil {
		log.Println("Error connecting to destination:", destErr)
		return destErr
	}
	defer destConn.Close()

	go copyMessages(conn, destConn)
	copyMessages(destConn, conn)
	return nil
}

func copyMessages(dst, src *websocket.Conn) {
	for {
		mt, message, err := src.ReadMessage()
		if err != nil {
			break
		}
		if err := dst.WriteMessage(mt, message); err != nil {
			break
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
