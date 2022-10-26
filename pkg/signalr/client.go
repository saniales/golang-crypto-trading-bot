package signalr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type negotiationResponse struct {
	Url                     string
	ConnectionToken         string
	ConnectionId            string
	KeepAliveTimeout        float32
	DisconnectTimeout       float32
	ConnectionTimeout       float32
	TryWebSockets           bool
	ProtocolVersion         string
	TransportConnectTimeout float32
	LogPollDelay            float32
}

type Client struct {
	OnMessageError func(err error)
	OnClientMethod func(hub, method string, arguments []json.RawMessage)
	// When client disconnects, the causing error is sent to this channel. Valid only after Connect().
	DisconnectedChannel chan bool
	params              negotiationResponse
	socket              *websocket.Conn
	nextId              int

	// Futures for server call responses and a guarding mutex.
	responseFutures map[string]chan *serverMessage
	mutex           sync.Mutex
	dispatchRunning bool
}

type serverMessage struct {
	Cursor     string            `json:"C"`
	Data       []json.RawMessage `json:"M"`
	Result     json.RawMessage   `json:"R"`
	Identifier string            `json:"I"`
	Error      string            `json:"E"`
}

func negotiate(scheme, address string) (negotiationResponse, error) {
	var response negotiationResponse

	var negotiationUrl = url.URL{Scheme: scheme, Host: address, Path: "/signalr/negotiate"}

	client := &http.Client{}

	reply, err := client.Get(negotiationUrl.String())
	if err != nil {
		return response, err
	}

	defer reply.Body.Close()

	if body, err := ioutil.ReadAll(reply.Body); err != nil {
		return response, err
	} else if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	} else {
		return response, nil
	}
}

func connectWebsocket(address string, params negotiationResponse, hubs []string) (*websocket.Conn, error) {
	var connectionData = make([]struct {
		Name string `json:"Name"`
	}, len(hubs))
	for i, h := range hubs {
		connectionData[i].Name = h
	}
	connectionDataBytes, err := json.Marshal(connectionData)
	if err != nil {
		return nil, err
	}

	var connectionParameters = url.Values{}
	connectionParameters.Set("transport", "webSockets")
	connectionParameters.Set("clientProtocol", "1.5")
	connectionParameters.Set("connectionToken", params.ConnectionToken)
	connectionParameters.Set("connectionData", string(connectionDataBytes))

	var connectionUrl = url.URL{Scheme: "wss", Host: address, Path: "signalr/connect"}
	connectionUrl.RawQuery = connectionParameters.Encode()

	if conn, _, err := websocket.DefaultDialer.Dial(connectionUrl.String(), nil); err != nil {
		return nil, err
	} else {
		return conn, nil
	}
}

func (self *Client) routeResponse(response *serverMessage) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if c, ok := self.responseFutures[response.Identifier]; ok {
		c <- response
		close(c)
		delete(self.responseFutures, response.Identifier)
	}
}

func (self *Client) createResponseFuture(identifier string) (chan *serverMessage, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if !self.dispatchRunning {
		return nil, fmt.Errorf("Dispatch is not running")
	}

	var c = make(chan *serverMessage)
	self.responseFutures[identifier] = c

	return c, nil
}

func (self *Client) deleteResponseFuture(identifier string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	delete(self.responseFutures, identifier)
}

func (self *Client) tryStartDispatch() error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.dispatchRunning {
		return fmt.Errorf("Another Dispatch() is running")
	}
	self.DisconnectedChannel = make(chan bool)
	self.dispatchRunning = true

	return nil
}

func (self *Client) endDispatch() {
	// Close all the waiting response futures.
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.dispatchRunning = false
	for _, c := range self.responseFutures {
		close(c)
	}
	self.responseFutures = make(map[string]chan *serverMessage)
	close(self.DisconnectedChannel)
}

// Start dispatch loop. This function will return when error occurs. When this
// happens, all the connections are closed and user can run Connect()
// and Dispatch() again on the same client.
func (self *Client) dispatch(connectedChannel chan bool) {
	if err := self.tryStartDispatch(); err != nil {
		panic("Dispatch is already running")
	}

	defer self.endDispatch()

	close(connectedChannel)

	for {
		var message serverMessage

		var hubCall struct {
			HubName   string            `json:"H"`
			Method    string            `json:"M"`
			Arguments []json.RawMessage `json:"A"`
		}

		_, data, err := self.socket.ReadMessage()
		if err != nil {
			self.socket.Close()
			break
		} else if err := json.Unmarshal(data, &message); err != nil {
			if self.OnMessageError != nil {
				self.OnMessageError(err)
			}
		} else {
			if len(message.Identifier) > 0 {
				// This is a response to a hub call.
				self.routeResponse(&message)
			} else if len(message.Data) == 1 {
				if err := json.Unmarshal(message.Data[0], &hubCall); err == nil && len(hubCall.HubName) > 0 && len(hubCall.Method) > 0 {
					// This is a client Hub method call from server.
					if self.OnClientMethod != nil {
						self.OnClientMethod(hubCall.HubName, hubCall.Method, hubCall.Arguments)
					}
				}
			}
		}
	}
}

// Call server hub method. Dispatch() function must be running, otherwise this method will never return.
func (self *Client) CallHub(hub, method string, params ...interface{}) (json.RawMessage, error) {
	var request = struct {
		Hub        string        `json:"H"`
		Method     string        `json:"M"`
		Arguments  []interface{} `json:"A"`
		Identifier int           `json:"I"`
	}{
		Hub:        hub,
		Method:     method,
		Arguments:  params,
		Identifier: self.nextId,
	}

	self.nextId++

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var responseKey = fmt.Sprintf("%d", request.Identifier)
	responseChannel, err := self.createResponseFuture(responseKey)
	if err != nil {
		return nil, err
	}

	if err := self.socket.WriteMessage(websocket.TextMessage, data); err != nil {
		return nil, err
	}

	defer self.deleteResponseFuture(responseKey)

	if response, ok := <-responseChannel; !ok {
		return nil, fmt.Errorf("Call to server returned no result")
	} else if len(response.Error) > 0 {
		return nil, fmt.Errorf("%s", response.Error)
	} else {
		return response.Result, nil
	}
}

func (self *Client) Connect(scheme, host string, hubs []string) error {
	// Negotiate parameters.
	if params, err := negotiate(scheme, host); err != nil {
		return err
	} else {
		self.params = params
	}

	// Connect Websocket.
	if ws, err := connectWebsocket(host, self.params, hubs); err != nil {
		return err
	} else {
		self.socket = ws
	}

	var connectedChannel = make(chan bool)
	go self.dispatch(connectedChannel)
	<-connectedChannel

	return nil
}

func (self *Client) Close() {
	self.socket.Close()
}

func NewWebsocketClient() *Client {
	return &Client{
		nextId:          1,
		responseFutures: make(map[string]chan *serverMessage),
	}
}
