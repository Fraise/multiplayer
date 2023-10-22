package main

import (
	"encoding/json"
	"log"
	"multiplayer/message"
	"multiplayer/player"
	"multiplayer/state"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func index(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("hello"))
}

func getWsHandler(state *state.State) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionId := vars["sessionId"]
		username := r.URL.Query().Get("username")

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer func() { _ = c.Close() }()

		p := &player.Player{
			Name:     username,
			Position: player.Position{},
		}
		err = state.AddPlayer(sessionId, p)

		if err != nil {
			log.Print("could not add player to the state\n", err)
			return
		}

		defer func() {
			err := state.RemovePlayer(sessionId, p)

			if err != nil {
				log.Printf("error removing player from state %s\n", err)
			}
		}()

		stopRead := make(chan bool, 1)
		stopWrite := make(chan bool, 1)

		// write loop
		go func() {
			defer log.Printf("stop sending data to %s\n", username)
			for {
				select {
				case <-time.After(time.Second):
					msg, err := state.GetSessionStateMessage(sessionId)

					if err != nil {
						log.Print("could not get session state message", err)
						stopRead <- true
						return
					}

					jsonMsg, err := json.Marshal(msg)

					if err != nil {
						log.Print("could not marshal session state message", err)
						continue
					}

					err = c.WriteMessage(websocket.TextMessage, jsonMsg)

					if err != nil {
						log.Print("error sending state: ", err)
						stopRead <- true
						return
					}
				case <-stopWrite:
					return
				}

			}
		}()

		defer log.Printf("stop reading data from %s\n", username)
		// read loop
		for {
			select {
			case <-stopRead:
				return
			default:
				_, msg, err := c.ReadMessage()

				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
						log.Printf("error reading player message: %v\n", err)
					}
					stopWrite <- true
					return
				}

				playerMsg := new(message.Player)
				err = json.Unmarshal(msg, playerMsg)

				if err != nil {
					log.Printf("error unmarshalling player message: %v\n", err)
					continue
				}

				p.Position.X = playerMsg.Position.X
				p.Position.Y = playerMsg.Position.Y
				p.Position.Z = playerMsg.Position.Z
			}

		}
	}
}

func main() {
	s := state.New(state.DefaultOptions())

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/sessions/{sessionId}", getWsHandler(s))
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
