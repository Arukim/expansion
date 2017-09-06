package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/arukim/expansion/game"
	"github.com/arukim/expansion/game/advisors"
	"github.com/arukim/expansion/models"
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

	advisors := []advisors.Advisor{
		advisors.NewExplorer(),
		advisors.NewGeneral(),
		advisors.NewInternal(),
	}

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
			b := game.NewBoard(&turnInfo)
			t := &models.Turn{
				Increase:  []models.Increase{},
				Movements: []models.Movement{},
			}
			for i, adv := range advisors {
				fmt.Printf("adv %v\n", i)
				time.Sleep(100 * time.Millisecond)
				adv.MakeTurn(b, t)
			}

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
