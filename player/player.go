package player

import (
	"fmt"
	"time"

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

	for i, adv := range p.advisors {
		fmt.Printf("adv %v\n", i)
		time.Sleep(100 * time.Millisecond)
		adv.MakeTurn(b, playerTurn)
	}

	return playerTurn
}
