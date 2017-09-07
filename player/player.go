package player

import (
	"fmt"

	"github.com/arukim/expansion/game"
	"github.com/arukim/expansion/game/advisors"
	"github.com/arukim/expansion/models"
)

type Player struct {
	id       int
	advisors []advisors.Advisor
}

// NewPlayer const
func NewPlayer(id int) *Player {
	p := &Player{id: id}

	p.advisors = []advisors.Advisor{
		advisors.NewExplorer(),
		advisors.NewGeneral(),
		advisors.NewInternal(),
	}

	return p
}

func (p *Player) MakeTurn(turnInfo *models.TurnInfo) *models.Turn {

	b := game.NewBoard(turnInfo)

	playerTurn := &models.Turn{
		Increase:  []models.Increase{},
		Movements: []models.Movement{},
	}

	fmt.Printf("P%d space: %d freeForces: %d\n", p.id, b.MyInfo.TerritorySize, b.MyInfo.ForcesFree)
	if b.MyInfo.ForcesTotal > 0 {
		for _, adv := range p.advisors {
			adv.MakeTurn(b, playerTurn)
		}
	} else {
		fmt.Printf("I've done")
	}
	fmt.Printf("P%d making turn\n", p.id)
	//fmt.Printf("player turn is %+v\n", playerTurn)
	return playerTurn
}
