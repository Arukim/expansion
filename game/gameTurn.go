package game

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ahmetb/go-linq"
	m "github.com/arukim/expansion/models"
)

type GameTurn struct {
	TurnInfo *m.TurnInfo
	Size     int
	Width    int

	WalkMap    *m.Map
	PlayersMap *m.Map
	MoveMap    *m.Map
	ForcesMap  *m.Map

	GoldList []m.Point
}

func NewGameTurn(t *m.TurnInfo) *GameTurn {
	g := new(GameTurn)

	g.parse(t)
	g.buildMoveMap()

	return g
}

func (g *GameTurn) FindMove() string {
	nearest := linq.From(g.GoldList).OrderByT(func(x m.Point) int {
		return g.GetDistance(x)
	}).First().(m.Point)

	return g.GetDirection(nearest)
}

func (g *GameTurn) parse(t *m.TurnInfo) {
	g.TurnInfo = t

	walkLayer := []rune(t.Layers[0])
	playersLayer := []rune(t.Layers[1])

	mapSize := len(walkLayer)
	mapWidth := int(math.Sqrt(float64(mapSize)))
	g.Size = mapSize
	g.Width = mapWidth

	fmt.Printf("Map size: %v, side: %v", mapSize, mapWidth)

	g.WalkMap = m.NewMap(mapWidth)
	g.GoldList = []m.Point{}

	for i := 0; i < mapSize; i++ {
		switch walkLayer[i] {
		case '$':
			g.GoldList = append(g.GoldList, m.NewPoint(i, g.Width))
			fallthrough
		case '1', '2', '3', '4', '.':
			g.WalkMap.Data[i] = 1
		}
	}

	g.PlayersMap = m.NewMap(mapWidth)
	for i := 0; i < mapSize; i++ {
		switch playersLayer[i] {
		case '♥':
			g.PlayersMap.Data[i] = 0
		case '♦':
			g.PlayersMap.Data[i] = 1
		case '♣':
			g.PlayersMap.Data[i] = 2
		case '♠':
			g.PlayersMap.Data[i] = 3
		case 45:
			g.PlayersMap.Data[i] = -1
		}
	}

	g.ForcesMap = m.NewMap(mapWidth)
	for i := 0; i < mapSize; i++ {
		value, _ := strconv.ParseInt(t.Forces[i*3:i*3+3], 36, 32)
		if value != 0 {
			g.ForcesMap.Data[i] = int(value)
		}
	}

	fmt.Println("Walk map")
	g.WalkMap.Print()
	fmt.Println("Players map")
	g.PlayersMap.Print()
	fmt.Println("Forces map")
	g.ForcesMap.Print()
}

func (g *GameTurn) buildMoveMap() {

	myForces := []m.Point{}

	g.PlayersMap.Iterate(func(i, v int) {
		if v == g.TurnInfo.MyColor {
			myForces = append(myForces, m.NewPoint(i, g.Width))
		}
	})

	fmt.Printf("My forces: %v\n", myForces)

	g.MoveMap = m.NewMap(g.Width)

	turn := 0
	for len(myForces) > 0 {
		turn++
		changes := []m.Point{}
		for _, f := range myForces {
			if g.MoveMap.Get(f) == 0 {
				g.MoveMap.Set(f, turn)
			}

			g.Neighbours(f, func(pos int, p m.Point) bool {
				moveV := g.MoveMap.Data[pos]
				walkV := g.WalkMap.Data[pos]

				if moveV == 0 && walkV == 1 {
					changes = append(changes, p)
				}

				return true
			})
		}
		myForces = changes
	}

	g.MoveMap.Print()
}

func (g *GameTurn) GetDistance(p m.Point) int {
	return g.WalkMap.Get(p) - 1
}

func (g *GameTurn) GetDirection(p m.Point) string {
	pos := g.MoveMap.Get(p)

	if pos == 0 {
		return "NONE"
	}

	if pos == 1 {
		return "STAY"
	}

	dir := ""
	for dir == "" {
		g.Neighbours(p, func(n_pos int, neighbour m.Point) bool {
			//found
			if g.MoveMap.Data[n_pos] == 1 {
				dir = neighbour.GetDirection(p)
				return false
			} else if g.MoveMap.Data[n_pos] < pos {
				p = neighbour
				pos = n_pos
				return false
			}
			return true
		})
	}

	return dir
}

func (g *GameTurn) Neighbours(p m.Point, f func(int, m.Point) bool) {
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
