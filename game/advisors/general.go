package advisors

import (
	linq "github.com/ahmetb/go-linq"
	"github.com/arukim/expansion/game"
	m "github.com/arukim/expansion/models"
)

// General contains game logic
type General struct {
	board *game.Board
}

func NewGeneral() *General {
	return &General{}
}

func (g *General) MakeTurn(b *game.Board, t *m.Turn) {
	g.board = b

	move := g.findMove()

	if move == nil {
		return
	}

	t.Increase = append(t.Increase, m.Increase{
		Count:  1000,
		Region: move.Region,
	})

	t.Movements = append(t.Movements, *move)
}

func (g *General) findMove() *m.Movement {
	var b = g.board

	p := m.Point{}

	if len(b.GoldList) > 0 {
		p = linq.From(b.GoldList).OrderByT(func(x m.Point) int {
			return b.GetDistance(x)
		}).First().(m.Point)
	} else if len(b.Enemies) > 0 {
		p = linq.From(b.Enemies).OrderByT(func(x m.Point) int {
			return b.GetDistance(x)
		}).First().(m.Point)
	} else {
		return nil
	}

	move := b.GetDirectionTo(p, b.OutsideMap)

	if move != nil {
		move.Count = b.ForcesMap.Get(move.Region) - 1
	}
	return move
}
