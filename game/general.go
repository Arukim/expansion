package game

import linq "github.com/ahmetb/go-linq"
import m "github.com/arukim/expansion/models"

// General contains game logic
type General struct {
	board *Board
}

func NewGeneral() *General {
	return &General{}
}

func (g *General) MakeTurn(b *Board) *m.Turn {
	g.board = b

	move := g.findMove()

	if move == nil {
		return nil
	}

	return &m.Turn{
		Increase: []m.Increase{
			m.Increase{
				Count:  10,
				Region: move.Region,
			},
		},
		Movements: []m.Movement{
			*move,
		},
	}
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
	}

	move := b.GetDirection(p)

	if move != nil {
		move.Count = b.ForcesMap.Get(move.Region) - 1
	}
	return move
}
