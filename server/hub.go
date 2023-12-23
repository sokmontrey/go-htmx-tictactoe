package server

import (
  "log"
  "net/http"
	"github.com/gorilla/websocket"
  "sokmontrey/go-htmx-tictactoe/game"
  "encoding/json"
  "strconv"
  "time"
  "fmt"
)

func generateUniqueID() string {
	return "guest_" + fmt.Sprint(time.Now().UnixNano())
}

type Player struct {
  conn *websocket.Conn
  hub *Hub
  send chan []byte
  id string
  mark string
}

func NewPlayer(h *Hub, c *websocket.Conn, id string) *Player {
  return &Player{
    conn: c,
    hub: h,
    send: make(chan []byte),
    id: id,
    mark: "X",
  }
}

func (p *Player) writePump() {
  for {
    select {
    case message := <- p.send:
      log.Println("Player: writePump: message " )
      p.conn.WriteMessage(websocket.TextMessage, message)
    }
  }
}

func (p *Player) readPump(){
  defer func(){
    p.conn.Close()
  }()
  for {
    _, message, err := p.conn.ReadMessage()
    if err != nil {
      // log.Printf("Player: %v", err)
      p.hub.leave_player <- p
      break
    }

    if p.hub.current_player != p {
      continue
    }

    if p.hub.player1 == nil || p.hub.player2 == nil {
      continue
    }

    if p.hub.player1 == p {
      p.hub.current_player = p.hub.player2
    } else {
      p.hub.current_player = p.hub.player1
    }

    p.hub.broadcast <- message
  }
}

type Hub struct {
  broadcast chan []byte
  player1 *Player
  player2 *Player
  new_player chan *Player
  leave_player chan *Player
  game *game.Game
  current_player *Player
}

func NewHub() *Hub {
  return &Hub{
    broadcast: make(chan []byte),
    player1: nil,
    player2: nil,
    new_player: make(chan *Player),
    leave_player: make(chan *Player),
    game: game.NewGame(),
    current_player: nil,
  }
}

func (h *Hub) run() {
  for {
    select {
    case new_player := <- h.new_player:
      log.Println("Hub: New Player " + new_player.id)
      if h.player1 == nil {
        new_player.mark = "X"
        h.player1 = new_player
      } else if h.player2 == nil {
        if h.player1.id == new_player.id {
          continue
        }
        new_player.mark = "O"
        h.player2 = new_player
      }
      h.current_player = new_player

    case leave_player := <- h.leave_player:
      log.Println("Hub: Player Disconnect " + leave_player.id)
      if h.player1 == leave_player {
        h.player1 = nil
      } else if h.player2 == leave_player {
        h.player2 = nil
      }

    case raw_data := <- h.broadcast:

      var data map[string]interface{}
      json.Unmarshal(raw_data, &data)
      square := data["square"].(string)
      square_num, _ := strconv.Atoi(square)
      is_end, game_message := h.game.Play(uint16(square_num), h.current_player.mark)

      log.Println("Hub: New Move: " + square)

      player_mark := "O"
      if is_end || game_message != "" {
        player_mark = game_message
      }

      symbol := game.Square(player_mark, square)
      if h.player1 != nil {
        select { case h.player1.send <- []byte(symbol): default: }
      }
      if h.player2 != nil {
        select { case h.player2.send <- []byte(symbol): default: }
      }
    }
  }
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(h *Hub, w http.ResponseWriter, r *http.Request){
  guestID, err := r.Cookie("guest_id")
	if err != nil || guestID.Value == "" {
		guestID = &http.Cookie{
			Name:    "guest_id",
			Value:   generateUniqueID(),
			Expires: time.Now().Add(30 * 24 * time.Hour), 
		}
		http.SetCookie(w, guestID)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
  
  player := NewPlayer(h, conn, guestID.Value)

  h.new_player <- player

  go player.writePump()
  go player.readPump()
}
