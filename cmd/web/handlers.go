package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const homeHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<title>Goction</title>
</head>
<body>
<div id="itemData">
<h1>{{.Name}}</h1>
<span id="price">{{.Price}}</span> {{.Currency}}
</div>
<script type="text/javascript">
(function() {
	var data = document.getElementById("price");
	var conn = new WebSocket("ws://{{.Host}}/socket");
	conn.onclose = function(evt) {
		console.log("Connection Closed");
	}
	conn.onmessage = function(evt) {
		console.log(evt.data);
		data.textContent = evt.data;
	}
})();
</script>
</body>
</html>
`

const (
	// Time allowed to write the price to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to update the price on site
	priceWait = time.Second
)

var (
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func errorBadRequest(w http.ResponseWriter, errLog *log.Logger, err error) {
	errLog.Println(err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host     string
		Name     string
		Price    uint64
		Currency string
	}{
		app.host,
		app.item.GetName(),
		app.item.GetPrice(),
		app.item.GetCurrency(),
	}
	homeTempl.Execute(w, &v)
}

func (app *application) socketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			app.errorLog.Println(err)
		}
		return
	}
	defer ws.Close()

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	go app.writer(ws)
	reader(ws)
}

func (app *application) updatePrice(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorBadRequest(w, app.errorLog, err)
		return
	}
	socketUrl := "ws://" + app.host + "/socket"

	price, err := strconv.ParseUint(string(reqBody), 10, 64)
	if err != nil {
		errorBadRequest(w, app.errorLog, err)
		return
	}

	err = app.item.UpdatePrice(uint64(price))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		app.errorLog.Println(err.Error())
		return
	}

	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		app.errorLog.Println("Error connecting to Websocket Server:", err)
		return
	}
	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, reqBody); err != nil {
		app.errorLog.Println(err)
		return
	}
}

// TODO(khatibomar): this should not be an app method!
func (app *application) writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	priceTicker := time.NewTicker(priceWait)
	defer func() {
		pingTicker.Stop()
		priceTicker.Stop()
		ws.Close()
	}()
	for {
		price := strconv.FormatUint(app.item.GetPrice(), 10)

		select {
		case <-priceTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.TextMessage, []byte(price)); err != nil {
				return
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
