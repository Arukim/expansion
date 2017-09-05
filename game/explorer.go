package game

import m "github.com/arukim/expansion/models"

// General contains game logic
type Explorer struct {
	board *Board
}

func NewExplorer() *Explorer {
	return &Explorer{}
}

func (g *Explorer) MakeTurn(b *Board, t *m.Turn) {

	for i := 0; i < b.Size; i++ {
		if b.MoveMap.Data[i] == 2 {
			p1 := m.NewPoint(i, b.Width)
			move := b.GetDirection(p1)
			move.Count = 1

			t.Movements = append(t.Movements, *move)
			t.Increase = append(t.Increase, m.Increase{
				Count:  1,
				Region: move.Region,
			})

		}
	}
}
