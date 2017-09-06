package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/arukim/expansion/models"
	"github.com/arukim/expansion/player"
	"github.com/gorilla/websocket"
)

type Client struct {
	Url string
}

func NewClient(url string) *Client {
	c := &Client{Url: url}

	go c.run()

	return c
}

func (c *Client) run() {

	u, err := url.Parse(c.Url)
	if err != nil {
		panic(err)
	}

	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	player := player.NewPlayer()
	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			turnInfo := models.TurnInfo{}
			json.Unmarshal(message[6:], &turnInfo)

			t := player.MakeTurn(&turnInfo)

			payload, _ := json.Marshal(t)
			msg := fmt.Sprintf("message('%s')", payload)
			log.Printf("%s\n", msg)
			//time.Sleep(100 * time.Millisecond)
			conn.WriteMessage(websocket.TextMessage, []byte(msg))
			//log.Printf("recv: %s", message)
			//log.Printf("turn info: %+v", turnInfo)
		}
	}()
}
