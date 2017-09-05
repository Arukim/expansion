package game

import m "github.com/arukim/expansion/models"

type Advisor interface {
	MakeTurn(b *Board, t *m.Turn)
}
