package player

import (
	"fmt"
	"log"

	"github.com/arukim/expansion/game"
	"github.com/arukim/expansion/game/advisors"
	"github.com/arukim/expansion/models"
)

type Player struct {
	id            int
	earlyAdvisors []advisors.Advisor
	lateAdvisors  []advisors.Advisor
}

// NewPlayer const
func NewPlayer(id int) *Player {
	p := &Player{id: id}

	p.earlyAdvisors = []advisors.Advisor{
		advisors.NewEarlyExplorer(),
		advisors.NewInternal(),
	}

	p.lateAdvisors = []advisors.Advisor{
		advisors.NewExplorer(),
		advisors.NewGeneral(),
		advisors.NewInternal(),
	}

	return p
}

func (p *Player) MakeTurn(turnInfo *models.TurnInfo) *models.Turn {

	b := game.NewBoard(turnInfo)

	if b.TotalWalkCells == b.MyInfo.TerritorySize {
		log.Println("won solo game")
		return nil
	}

	playerTurn := &models.Turn{
		Increase:  []models.Increase{},
		Movements: []models.Movement{},
	}

	fmt.Printf("P%d inc: %d space: %d freeForces: %d occup: %f\n", p.id, b.ForcesAvailable, b.MyInfo.TerritorySize, b.MyInfo.ForcesFree, b.OccupationRate)
	if b.MyInfo.ForcesTotal > 0 {
		var advSet []advisors.Advisor
		// TODO: add more checks, enclave territory check, free mines check.
		if b.OccupationRate < 0.75 {
			advSet = p.earlyAdvisors
		} else {
			advSet = p.lateAdvisors
		}

		for _, adv := range advSet {
			adv.MakeTurn(b, playerTurn)
		}

	} else {
		fmt.Printf("I've done")
	}
	//fmt.Printf("P%d making turn\n", p.id)
	//fmt.Printf("player turn is %+v\n", playerTurn)
	return playerTurn
}
