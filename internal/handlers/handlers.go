package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response sent back from websocket server
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndpoint upgrade connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Connection cannot be upgrade connection\n" + err.Error())
	}

	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Message = "<em><small>Connected to server</small></em>"
	
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""
	
	response.Action = "Connect"
	response.ConnectedUsers = getUserList()

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("Connection cannot be upgrade connection\n" + err.Error())
	}

	go ListenForWsConnection(&conn)
}

func ListenForWsConnection(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("err", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		e := <-wsChan

		switch e.Action {
			case "username":
				clients[e.Conn] = e.Username
				var users []string = getUserList()
				response.Action = "list_users"
				response.ConnectedUsers = users
				broadcastToAll(response)

			case "left":
				delete(clients, e.Conn)
				response.Action = "list_users"
				var users []string = getUserList()
				response.ConnectedUsers = users
				broadcastToAll(response)

			case "broadcast":
				response.Action = "broadcast"
				response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
				broadcastToAll(response)
		}

	}
}

func getUserList() []string {
	var userList []string

	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}

	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("websocket error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
