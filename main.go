package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/arukim/expansion/game"
	"github.com/arukim/expansion/game/advisors"
	"github.com/arukim/expansion/models"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//u, err := url.Parse("ws://127.0.0.1:8080/codenjoy-contest/ws?user=1@a.com")
	u, err := url.Parse("ws://ecsc00104eef.epam.com:8080/codenjoy-contest/ws?user=nikita_smelov2@epam.com")
	if err != nil {
		panic(err)
	}

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	advisors := []advisors.Advisor{
		advisors.NewExplorer(),
		advisors.NewGeneral(),
		advisors.NewInternal(),
	}

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
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
			for _, adv := range advisors {
				adv.MakeTurn(b, t)
			}

			payload, _ := json.Marshal(t)
			msg := fmt.Sprintf("message('%s')", payload)
			log.Printf("%s\n", msg)
			//time.Sleep(100 * time.Millisecond)
			c.WriteMessage(websocket.TextMessage, []byte(msg))
			//log.Printf("recv: %s", message)
			//log.Printf("turn info: %+v", turnInfo)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
