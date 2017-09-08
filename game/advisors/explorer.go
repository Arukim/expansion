package advisors

import (
	"math/rand"

	linq "github.com/ahmetb/go-linq"
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
	inc := make(map[m.Point]bool)
	for i := 0; i < b.Size; i++ {
		if b.OutsideMap.Data[i] == 2 {
			p1 := m.NewPoint(i, b.Width)
			targetForces := b.ForcesMap.Get(p1)
			moves := b.GetDirectionTo(p1, b.OutsideMap)

			for _, mov := range moves {
				inc[mov.Region] = true
			}

			movesL := linq.From(moves)
			// check if zerg atack is available
			myTotalForces := int(movesL.SelectT(func(m m.Movement) int {
				return b.ForcesMap.Get(m.Region) - 1
			}).SumInts())

			if myTotalForces > targetForces {
				total := targetForces
				for _, m := range moves {
					availableToAttack := b.ForcesMap.Get(m.Region) - 1

					if availableToAttack == 0 {
						continue
					}

					if availableToAttack > total {
						rest := availableToAttack - total
						if rest > 2 && rest < 20 {
							rest = 2
						}
						m.Count = total + rest
					} else {
						m.Count = availableToAttack
					}

					total -= m.Count
					b.ForcesMap.Data[m.Region.GetPos(b.Width)] -= m.Count

					t.Movements = append(t.Movements, m)

					if total <= 0 {
						break
					}
				}
			}
		}
	}

	dest := make([]m.Increase, len(inc))
	perm := rand.Perm(len(inc))
	i := 0
	for p := range inc {
		dest[perm[i]] = m.Increase{Count: 1, Region: p}
		i++
	}

	t.Increase = append(t.Increase, dest...)
}
