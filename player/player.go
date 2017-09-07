package player

import (
	"fmt"

	"github.com/arukim/expansion/game"
	"github.com/arukim/expansion/game/advisors"
	"github.com/arukim/expansion/models"
)

type Player struct {
	advisors []advisors.Advisor
}

// NewPlayer const
func NewPlayer() *Player {
	p := &Player{}

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

	for _, adv := range p.advisors {
		adv.MakeTurn(b, playerTurn)
	}

	fmt.Printf("player turn is %+v\n", playerTurn)
	return playerTurn
}
