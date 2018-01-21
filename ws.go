package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Command struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

func coinpkimain(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, jsonstr, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", jsonstr)
		var command Command
		json.Unmarshal(jsonstr, &command)
		fmt.Println(command)
		if command.Method == "prove" && len(command.Params) == 3 {
			walletaddr := command.Params[0]
			prooftext := command.Params[1]
			sign := command.Params[2]
			find(walletaddr, sign, prooftext)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/coinpkiws", coinpkimain)
	addr := flag.String("addr", "localhost:8088", "http service address")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
