package advisors

import (
	"github.com/arukim/expansion/game"
	m "github.com/arukim/expansion/models"
)

// General contains game logic
type Explorer struct {
}

func NewExplorer() *Explorer {
	return &Explorer{}
}

func (g *Explorer) MakeTurn(b *game.Board, t *m.Turn) {

	for i := 0; i < b.Size; i++ {
		if b.OutsideMap.Data[i] == 2 {
			p1 := m.NewPoint(i, b.Width)
			move := b.GetDirectionTo(p1, b.OutsideMap)
			move.Count = 1

			t.Movements = append(t.Movements, *move)
			t.Increase = append(t.Increase, m.Increase{
				Count:  1,
				Region: move.Region,
			})

		}
	}
}
