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

	moves := g.findMove()

	if len(moves) == 0 {
		return
	}

	t.Increase = append(t.Increase, m.Increase{
		Count:  1000,
		Region: moves[0].Region,
	})

	t.Movements = append(t.Movements, moves...)
}

func (g *General) findMove() []m.Movement {
	var b = g.board

	p := m.Point{}

	if len(b.GoldList) > 0 {
		p = linq.From(b.GoldList).OrderByT(func(x m.Point) int {
			return b.OutsideMap.Get(x) - 1
		}).First().(m.Point)
	} else if len(b.Enemies) > 0 {
		p = linq.From(b.Enemies).OrderByT(func(x m.Point) int {
			return b.OutsideMap.Get(x) - 1
		}).First().(m.Point)
	} else {
		return []m.Movement{}
	}

	moves := b.GetDirectionTo(p, b.OutsideMap)

	// don't forget to move largest force!
	maxForce := 0
	maxForcePos := 0
	for i, v := range b.ForcesMap.Data {
		// check only my force
		if b.PlayersMap.Data[i] == b.TurnInfo.MyColor && v > maxForce {
			maxForce = v
			maxForcePos = i
		}
	}

	moves = append(moves, *b.GetDirectionFromTo(m.NewPoint(maxForcePos, b.Width), p))

	for i := range moves {
		moves[i].Count = b.ForcesMap.Get(moves[i].Region) - 1
	}

	//fmt.Printf("General moves: %+v\n", moves)
	return moves
}
