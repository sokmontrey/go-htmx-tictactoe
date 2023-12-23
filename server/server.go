package server

import (
  "sokmontrey/go-htmx-tictactoe/game"
  "net/http"
  "html/template"
  "log"
  "strconv"
)

func homePage(w http.ResponseWriter, r *http.Request){
  t := template.Must(template.ParseFiles("./templates/index.html")) 
  t.Execute(w, nil)
}

func gameBoard(w http.ResponseWriter, r *http.Request){
  board := ""
  for i := 0; i < 9; i++ {
    board += game.Square(" ", strconv.Itoa(i)) 
  }
  tmpl, _ := template.New("board").Parse(board)
  tmpl.Execute(w, nil)
}

func StartServer() {
  hub := NewHub()
  go hub.run()

  http.HandleFunc("/", homePage)

  http.HandleFunc("/board", gameBoard)

  http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request){
    wsHandler(hub, w, r)
  })

  log.Println("App is running on :8080...")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
