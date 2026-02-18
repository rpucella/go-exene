package exene

import (
	"time"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"github.com/gorilla/websocket"
	"encoding/json"
)


/*
   ************************************************************
   
     Go eXene library

   ************************************************************
*/


type App struct {
	Url string
	listener net.Listener
	mux http.Handler
	browser []string
}

func NewBrowserApp(d Dispatcher, w Widget, browser []string) *App {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)
	log.Printf("Listening on port %d\n", port)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", GetPage(port))
	mux.HandleFunc("GET /sdk.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/sdk.js") })
	mux.HandleFunc("GET /socket", WebSocketHandler(d, w))
	return &App{url, listener, mux, browser}
}

func (app *App) Start() {
	go func(){
		browser := app.browser
		delay := time.Duration(2)
		time.Sleep(delay * time.Second)
		log.Printf("Executing %s\n", browser[0])
		browser = append(browser, app.Url)
		cmd := exec.Command(browser[0], browser[1:]...)
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}()
	panic(http.Serve(app.listener, app.mux))
}

func GetPage(port int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`
<!doctype html>
<html>
  <head>
   <script src="/sdk.js"></script>
   <script>window.onload = () => { connect("ws://localhost:%d/socket") }</script>
  </head>
  <body style="margin: 0;">
  </body>
</html>
`, port)))
	}
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func WebSocketHandler(d Dispatcher, widget Widget) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Websocket upgrade request %s\n", r.URL)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		conn.SetCloseHandler(func(code int, text string) error {
			writeWait := time.Second
			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""), time.Now().Add(writeWait))
			return nil
		})
		outgoing2 := struct{Type string `json:"type"`; Widget Widget `json:"widget"`}{"widget", widget}
		msg2, err := json.Marshal(outgoing2)
		if err != nil {
			return
		}
		conn.WriteMessage(websocket.TextMessage, msg2)
		go MessagePump(conn, d)
		for {
			select {
			case obj := <- d.GetEvent():
				if obj["type"].(string) == "event" {
				}
				// Assume click event.
				target := obj["target"].(string)
				d.DispatchEvent(target)
			case obj := <- d.PutUpdate():
				outgoing, err := json.Marshal(obj)
				if err != nil {
					// Ignore.
					continue
				}
				conn.WriteMessage(websocket.TextMessage, outgoing)
			}
		}
	}
}

type Dispatcher interface {
	GetEvent () chan map[string]any
	PutUpdate () chan map[string]any
	DispatchEvent (id string)
	RegisterEvent (id string, event chan bool)
}

type Dispatch struct {
	EventChan chan map[string]any
	UpdateChan chan map[string]any
	Table map[string]chan bool
}

func MessagePump(conn *websocket.Conn, d Dispatcher) {
	// Message pump.
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		msg := make(map[string]any)
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("read:", err)
			continue
		}
		d.GetEvent() <- msg
	}
}

func NewDispatcher() *Dispatch {
	event := make(chan map[string]any)
	update := make(chan map[string]any)
	table := make(map[string]chan bool)
	return &Dispatch{event, update, table}
}

func (d *Dispatch) GetEvent() chan map[string]any {
	return d.EventChan
}

func (d *Dispatch) PutUpdate() chan map[string]any {
	return d.UpdateChan
}

func (d *Dispatch) DispatchEvent(id string) {
	d.Table[id] <- true
}

func (d *Dispatch) RegisterEvent(id string, event chan bool) {
	d.Table[id] = event
}
