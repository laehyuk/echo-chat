package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var upgrader = websocket.Upgrader{}	//Convert Http to Websocket

func main()  {
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func echo(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print("upgrade : ", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err!= nil {
			log.Printf("read : ", err)
			break
		}
		var objmap map[string]interface{}
		_ = json.Unmarshal(message, &objmap)	// decode
		event := objmap["event"].(string)
		sendData := map[string]interface{}{
			"event": "res",
			"data": nil,
		}

		switch event {
		case "open":
			log.Printf("Received : %s\n", event)
		case "req":
			sendData["data"] = objmap["data"]
			log.Printf("Received : %s\n", event)
		}

		refineSendData, err := json.Marshal(sendData)	// encode
		err = conn.WriteMessage(messageType, refineSendData)
		if err != nil {
			log.Printf("write : ", err)
			break
		}
	}
}

func index(writer http.ResponseWriter, request *http.Request)  {
	path := filepath.Join("templates", "index.html")
	template := template.Must(template.ParseFiles(path))
	template.Execute(writer, "ws://"+request.Host+"/echo")
}
