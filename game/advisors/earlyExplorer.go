package advisors

import (
	"log"

	"github.com/arukim/expansion/game"
	m "github.com/arukim/expansion/models"
)

// EarlyExplorer is advisor with main focus to capture new territories
type EarlyExplorer struct {
}

func NewEarlyExplorer() *EarlyExplorer {
	return &EarlyExplorer{}
}

func (g *EarlyExplorer) MakeTurn(b *game.Board, t *m.Turn) {
	// build reversed-flood fill
	// clone outside map
	reverseMap := b.OutsideMap.Clone()

	pointsMap := make(map[m.Point]int)
	// remove all enemies from map
	reverseMap.Modify(func(p m.Point, v int) int {
		if _, ok := b.Enemies[p]; ok {
			return -1
		}
		return v
	})

	max := 0
	// for every free cell add 1 point into pointsMap
	reverseMap.IterateP(func(p m.Point, v int) {
		if v > 0 {
			pointsMap[p] = 1
		}

		if v > max {
			max = v
		}
	})

	// reverse flood fill, starting from max point
	for i := max; i > 0; i-- {
		reverseMap.IterateP(func(p m.Point, v int) {
			if v > 0 {
				b.Neighbours(p, func(_ int, p1 m.Point) bool {
					log.Println("inc")
					if reverseMap.Get(p1) < i {
						pointsMap[p1] += pointsMap[p]
						log.Println("inc")
						return false
					}
					return true
				})
			}
		})
	}

	// reverse flood-fill

	// drop all available forces
	// do a turn
}
