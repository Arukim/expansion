package game

import (
	"fmt"
	"math"
	"strconv"

	m "github.com/arukim/expansion/models"
)

type Board struct {
	TurnInfo *m.TurnInfo
	Size     int
	Width    int

	WalkMap    *m.Map
	PlayersMap *m.Map
	MoveMap    *m.Map
	ForcesMap  *m.Map

	GoldList []m.Point
}

func NewBoard(t *m.TurnInfo) *Board {
	b := new(Board)

	b.parse(t)
	b.buildMoveMap()

	return b
}

func (b *Board) parse(t *m.TurnInfo) {
	b.TurnInfo = t

	walkLayer := []rune(t.Layers[0])
	playersLayer := []rune(t.Layers[1])

	mapSize := len(walkLayer)
	mapWidth := int(math.Sqrt(float64(mapSize)))
	b.Size = mapSize
	b.Width = mapWidth

	fmt.Printf("Map size: %v, side: %v", mapSize, mapWidth)

	b.WalkMap = m.NewMap(mapWidth)
	b.GoldList = []m.Point{}

	for i := 0; i < mapSize; i++ {
		switch walkLayer[i] {
		case '$':
			b.GoldList = append(b.GoldList, m.NewPoint(i, b.Width))
			fallthrough
		case '1', '2', '3', '4', '.':
			b.WalkMap.Data[i] = 1
		}
	}

	b.PlayersMap = m.NewMap(mapWidth)
	for i := 0; i < mapSize; i++ {
		switch playersLayer[i] {
		case '♥':
			b.PlayersMap.Data[i] = 0
		case '♦':
			b.PlayersMap.Data[i] = 1
		case '♣':
			b.PlayersMap.Data[i] = 2
		case '♠':
			b.PlayersMap.Data[i] = 3
		case 45:
			b.PlayersMap.Data[i] = -1
		}
	}

	b.ForcesMap = m.NewMap(mapWidth)
	for i := 0; i < mapSize; i++ {
		value, _ := strconv.ParseInt(t.Forces[i*3:i*3+3], 36, 32)
		if value != 0 {
			b.ForcesMap.Data[i] = int(value)
		}
	}

	fmt.Println("Walk map")
	b.WalkMap.Print()
	fmt.Println("Players map")
	b.PlayersMap.Print()
	fmt.Println("Forces map")
	b.ForcesMap.Print()
}

func (b *Board) buildMoveMap() {

	myForces := []m.Point{}

	b.PlayersMap.Iterate(func(i, v int) {
		if v == b.TurnInfo.MyColor {
			myForces = append(myForces, m.NewPoint(i, b.Width))
		}
	})

	fmt.Printf("My forces: %v\n", myForces)

	b.MoveMap = m.NewMap(b.Width)

	turn := 0
	for len(myForces) > 0 {
		turn++
		changes := []m.Point{}
		for _, f := range myForces {
			if b.MoveMap.Get(f) == 0 {
				b.MoveMap.Set(f, turn)
			}

			b.Neighbours(f, func(pos int, p m.Point) bool {
				moveV := b.MoveMap.Data[pos]
				walkV := b.WalkMap.Data[pos]

				if moveV == 0 && walkV == 1 {
					changes = append(changes, p)
				}

				return true
			})
		}
		myForces = changes
	}

	b.MoveMap.Print()
}

func (b *Board) GetDistance(p m.Point) int {
	return b.WalkMap.Get(p) - 1
}

func (b *Board) GetDirection(p m.Point) *m.Movement {
	pos := b.MoveMap.Get(p)

	if pos == 0 {
		return nil
	}

	if pos == 1 {
		return nil
	}

	dir := ""
	for dir == "" {
		b.Neighbours(p, func(n_pos int, neighbour m.Point) bool {
			//found
			if b.MoveMap.Data[n_pos] == 1 {
				dir = neighbour.GetDirection(p)
				p = neighbour
				return false
			} else if b.MoveMap.Data[n_pos] == pos-1 {
				p = neighbour
				pos--
				return false
			}
			return true
		})
	}

	return &m.Movement{
		Direction: dir,
		Region:    p,
	}
}

func (b *Board) Neighbours(p m.Point, f func(int, m.Point) bool) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			p1 := p.Add(i, j)
			if p1.X < 0 || p1.Y < 0 || p1.X >= b.Width || p1.Y >= b.Width {
				continue
			}
			pos := p1.GetPos(b.Width)
			if !f(pos, p1) {
				return
			}
		}
	}
}
