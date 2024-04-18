// Copyright 2024 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/livekit/protocol/livekit"
	"github.com/livekit/signal-proxy/pkg/config"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	writer                http.ResponseWriter
	request               *http.Request
	destinationLiveKitURL string
	rewriteIceServers     []*livekit.ICEServer
}

func NewConnection(
	config *config.Config,
	writer http.ResponseWriter,
	request *http.Request,
) (*Connection, error) {

	newIceServers := make([]*livekit.ICEServer, 0)
	for _, iceServer := range config.ICEServers {
		newIceServers = append(newIceServers, &livekit.ICEServer{
			Urls:       iceServer.Urls,
			Username:   iceServer.Username,
			Credential: iceServer.Credential,
		})
	}

	return &Connection{
		writer:                writer,
		request:               request,
		destinationLiveKitURL: config.DestinationLiveKitURL,
		rewriteIceServers:     newIceServers,
	}, nil
}

func (c *Connection) Run() error {
	conn, err := upgrader.Upgrade(c.writer, c.request, nil)
	if err != nil {
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
	parsed, err := url.Parse(c.destinationLiveKitURL)
	if err != nil {
		return fmt.Errorf("error parsing destination URL: %w", err)
	}
	host := parsed.Host
	scheme := parsed.Scheme
	destURL := url.URL{Scheme: scheme, Host: host, Path: c.request.URL.Path, RawQuery: queryParams}

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
		return fmt.Errorf("error connecting to destination: %w", destErr)
	}
	defer destConn.Close()

	go c.copyMessages(destConn, conn)
	c.copyServerMessages(conn, destConn)
	return nil
}

func (c *Connection) copyServerMessages(dst, src *websocket.Conn) {
	for {
		mt, message, err := src.ReadMessage()

		if err != nil {
			break
		}

		newMessage, err := c.modifyServerMessage(message)
		if err != nil {
			log.Printf("Error modifying message: %v", err)
			break
		}

		if err := dst.WriteMessage(mt, newMessage); err != nil {
			log.Printf("Error writing message to destination: %v", err)
			break
		}
	}
}

func (c *Connection) modifyServerMessage(msg []byte) ([]byte, error) {
	signalResponse := &livekit.SignalResponse{}
	err := proto.Unmarshal(msg, signalResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal SignalResponse message: %w", err)
	}

	updated := false
	if join := signalResponse.GetJoin(); join != nil {
		join.IceServers = c.rewriteIceServers
		updated = true
	} else if reconnect := signalResponse.GetReconnect(); reconnect != nil {
		reconnect.IceServers = c.rewriteIceServers
		updated = true
	}

	// Save some work if we didn't update anything
	if !updated {
		return msg, nil
	}

	modifiedMessage, err := proto.Marshal(signalResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reconnect SignalResponse message: %w", err)
	}
	return modifiedMessage, nil
}

func (c *Connection) copyMessages(dst, src *websocket.Conn) {
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
