package game

import (
	"fmt"
	"math"

	"github.com/ahmetb/go-linq"
	"github.com/arukim/expansion/models"
)

type GameTurn struct {
	TurnInfo *models.TurnInfo
	Size     int
	Width    int

	WalkMap   []int
	ForcesMap []int
	MoveMap   []int

	GoldMap []models.Point
}

func NewGameTurn(t *models.TurnInfo) *GameTurn {
	g := new(GameTurn)

	g.parse(t)
	g.buildMoveMap()

	return g
}

func (g *GameTurn) FindMove() string {
	nearest := linq.From(g.GoldMap).OrderByT(func(x models.Point) int {
		return g.GetDistance(x)
	}).First().(models.Point)

	return g.GetDirection(nearest)
}

func (g *GameTurn) parse(t *models.TurnInfo) {
	g.TurnInfo = t

	walkLayer := []rune(t.Layers[0])
	forceLayer := []rune(t.Layers[1])

	mapSize := len(walkLayer)
	mapWidth := int(math.Sqrt(float64(mapSize)))
	g.Size = mapSize
	g.Width = mapWidth

	fmt.Printf("Map size: %v, side: %v", mapSize, mapWidth)

	g.WalkMap = make([]int, mapSize)
	g.GoldMap = []models.Point{}

	for i := 0; i < mapSize; i++ {
		switch walkLayer[i] {
		case '$':
			g.GoldMap = append(g.GoldMap, models.NewPoint(i, g.Width))
			fallthrough
		case '1', '2', '3', '4', '.':
			g.WalkMap[i] = 1

		}
	}

	g.ForcesMap = make([]int, mapSize)
	for i := 0; i < mapSize; i++ {
		switch forceLayer[i] {
		case '♥':
			g.ForcesMap[i] = 0
		case '♦':
			g.ForcesMap[i] = 1
		case '♣':
			g.ForcesMap[i] = 2
		case '♠':
			g.ForcesMap[i] = 3
		default:
			g.ForcesMap[i] = -1
		}
	}

	for i := 0; i < mapWidth; i++ {
		fmt.Printf("%v\n", g.WalkMap[i*mapWidth:i*mapWidth+mapWidth])
	}

	for i := 0; i < mapWidth; i++ {
		fmt.Printf("%v\n", g.ForcesMap[i*mapWidth:i*mapWidth+mapWidth])
	}

}

func (g *GameTurn) buildMoveMap() {

	myForces := []models.Point{}

	for i := 0; i < g.Size; i++ {
		if g.ForcesMap[i] == g.TurnInfo.MyColor {
			myForces = append(myForces, models.Point{X: i % g.Width, Y: i / g.Width})
		}
	}

	fmt.Printf("My forces: %v\n", myForces)

	g.MoveMap = make([]int, g.Size)
	turn := 0
	for len(myForces) > 0 {
		turn++
		changes := []models.Point{}
		for _, f := range myForces {
			g.MoveMap[f.GetPos(g.Width)] = turn
			g.Neighbours(f, func(pos int, p models.Point) bool {
				moveV := g.MoveMap[pos]
				walkV := g.WalkMap[pos]

				if moveV == 0 && walkV == 1 {
					changes = append(changes, p)
				}

				return true
			})
		}
		myForces = changes
	}

	for i := 0; i < g.Width; i++ {
		fmt.Printf("%v\n", g.MoveMap[i*g.Width:i*g.Width+g.Width])
	}
}

func (g *GameTurn) GetDistance(p models.Point) int {
	return g.WalkMap[p.GetPos(g.Width)] - 1
}

func (g *GameTurn) GetDirection(p models.Point) string {
	pos := g.MoveMap[p.GetPos(g.Width)]

	if pos == 0 {
		return "NONE"
	}

	if pos == 1 {
		return "STAY"
	}

	dir := ""
	for dir == "" {
		g.Neighbours(p, func(n_pos int, neighbour models.Point) bool {
			//found
			if g.MoveMap[n_pos] == 1 {
				dir = neighbour.GetDirection(p)
				return false
			} else if g.MoveMap[n_pos] < pos {
				p = neighbour
				pos = n_pos
				return false
			}
			return true
		})
	}
	return dir
}

func (g *GameTurn) Neighbours(p models.Point, f func(int, models.Point) bool) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			p1 := p.Add(i, j)
			if p1.X < 0 || p1.Y < 0 || p1.X >= g.Width || p1.Y >= g.Width {
				continue
			}
			pos := p1.GetPos(g.Width)
			if !f(pos, p1) {
				return
			}

		}
	}
}
