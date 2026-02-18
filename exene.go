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

func NewBrowserApp(sh Shell, browser []string) *App {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)
	log.Printf("listening on port %d\n", port)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", GetPage(port))
	mux.HandleFunc("GET /sdk.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/sdk.js") })
	mux.HandleFunc("GET /socket", WebSocketHandler(sh))
	return &App{url, listener, mux, browser}
}

func (app *App) Start() {
	go func(){
		browser := app.browser
		delay := time.Duration(2)
		time.Sleep(delay * time.Second)
		log.Printf("executing %s\n", browser[0])
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
   <style>.debug { box-sizing: border-box; border: 1px solid red; }</style>
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

func WebSocketHandler(sh Shell) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("websocket upgrade request %s\n", r.URL)
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
		eventChan := make(chan map[string]any)
		updateChan := make(chan map[string]any)
		dispatchMap := make(map[string]chan bool)
		webIfc := &WebInterface{eventChan, updateChan, dispatchMap}
		// Read size off the first message coming from the frontend advertising the size!
		// Also create the environmental channels
		_, initMessage, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		initMsg := make(map[string]any)
		err = json.Unmarshal(initMessage, &initMsg)
		if err != nil {
			log.Println("read:", err)
			return
		}
		if initMsg["type"] != "init" {
			log.Println("init message type:", initMsg["type"])
			return
		}
		width := int(initMsg["width"].(float64))
		height := int(initMsg["height"].(float64))
		log.Printf("viewport size %d x %d\n", width, height)
		html := sh.Init(webIfc, width, height)
		outgoing2 := struct{Type string `json:"type"`; Widget Html `json:"widget"`}{"widget", html}
		msg2, err := json.Marshal(outgoing2)
		if err != nil {
			return
		}
		conn.WriteMessage(websocket.TextMessage, msg2)
		go WebSocketMessagePump(conn, webIfc)
		for {
			select {
			case obj := <- eventChan:
				if obj["type"].(string) == "event" {
				}
				// Assume click event.
				target := obj["target"].(string)
				// Dispatch by sending to the shell instead?
				dispatchMap[target] <- true
				
			case obj := <- updateChan:
				outgoing, err := json.Marshal(obj)
				if err != nil {
					// Ignore.
					continue
				}
				// Updates = change label/text
				// Updates = change size
				// Batch updates!
				conn.WriteMessage(websocket.TextMessage, outgoing)
			}
		}
	}
}

type WebInterface struct {
	eventChan chan map[string]any
	updateChan chan map[string]any
	dispatchMap map[string]chan bool
}

func WebSocketMessagePump(conn *websocket.Conn, webIfc *WebInterface) {
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
		webIfc.eventChan <- msg
	}
}
