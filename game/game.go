package game

import (
  "fmt"
)

func Square(char string, id string) string {
  return `<button id="square`+id+`" name="square" value="`+id+`" type="submit"
    class="p-5 text-[2rem] bg-[white] opacity-50 hover:opacity-100"
  >`+char+`</button>`
}

type Board struct {
  o uint16
  x uint16
} 

func (b *Board) isOccupied(square uint16) bool {
  pos := uint16(1) << square
  return (b.o & pos != 0) || (b.x & pos != 0)
}

func (b *Board) isFull() bool {
  return (b.o | b.x) == 0x1FF
}

func (b *Board) checkWinner(is_o_turn bool) bool {
  current_player := b.x
  if is_o_turn { current_player = b.o }

  // check rows
  for i:=0; i<3; i++ {
    comparer := uint16(0b111) << (i*3)
    if (current_player & comparer) == comparer { 
      return true 
    } 
  }

  // check cols
  for i:=0; i<3; i++ {
    comparer := uint16(0b1001001) << (i)
    if (current_player & comparer) == comparer {
      return true
    }
  }

  //check diagonal
  if (current_player & uint16(0b100010001) == uint16(273)) || (current_player & uint16(0b1010100) == uint16(84)) {
    return true
  }

  return false
}

func (b *Board) placeMark(square uint16, is_o_turn bool) {
  pos := uint16(1) << square
  if is_o_turn {
    b.o |= pos
  } else {
    b.x |= pos
  }
}

type Game struct {
  board Board
}

func NewGame() *Game {
  return &Game{
    board: Board{ o: 0, x: 0 }, 
  }
}

func (g *Game) Play(square uint16, player_mark string) (bool, string) {
  if g.board.isOccupied(square) {
    return false, player_mark
  }

  is_o_turn := player_mark == "O"

  g.board.placeMark(square, is_o_turn)

  if g.board.checkWinner(is_o_turn) {
    if is_o_turn { 
      return true, "O is the winner!" 
    }
    return true, "X is the winner!" 
  }

  if g.board.isFull(){
    return true, "Draw!"
  }

  return false, player_mark
}

func (g *Game) Show() {
  b := g.board
  for i:=0; i<3; i++ {
    for j:=0; j<3; j++ {
      square := i*3 + j
      pos := uint16(1) << square

      var player string
      if b.o & pos != 0 {
        player = "O"
      } else if b.x & pos != 0 {
        player = "X"
      } else {
        player = "_"
      }

      fmt.Print(player)
    }
    fmt.Println()
  }
}
