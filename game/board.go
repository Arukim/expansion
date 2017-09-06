package game

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ahmetb/go-linq"

	m "github.com/arukim/expansion/models"
)

type Board struct {
	TurnInfo *m.TurnInfo
	Size     int
	Width    int

	WalkMap    *m.Map
	PlayersMap *m.Map
	OutsideMap *m.Map
	InsideMap  *m.Map
	ForcesMap  *m.Map

	GoldList []m.Point
	Enemies  []m.Point
}

// NewBoard instance creation
func NewBoard(t *m.TurnInfo) *Board {
	b := new(Board)

	b.parse(t)
	b.buildOutsideMap()
	b.buildInsideMap()

	return b
}

func (b *Board) rotate(i int) int {
	return (i % b.Width) + (b.Width-1-i/b.Width)*b.Width
}

func (b *Board) parse(t *m.TurnInfo) {
	b.TurnInfo = t

	walkLayer := []rune(t.Layers[0])
	playersLayer := []rune(t.Layers[1])

	mapSize := len(walkLayer)
	mapWidth := int(math.Sqrt(float64(mapSize)))

	b.Size = mapSize
	b.Width = mapWidth

	//fmt.Printf("Map size: %v, side: %v", mapSize, mapWidth)

	b.WalkMap = m.NewMap(mapWidth)
	b.GoldList = []m.Point{}

	for i := 0; i < mapSize; i++ {
		switch walkLayer[b.rotate(i)] {
		case '$':
			b.GoldList = append(b.GoldList, m.NewPoint(i, b.Width))
			fallthrough
		case '1', '2', '3', '4', '.':
			b.WalkMap.Data[i] = 1
		}
	}

	b.PlayersMap = m.NewMap(mapWidth)
	for i := 0; i < mapSize; i++ {
		switch playersLayer[b.rotate(i)] {
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
		if b.PlayersMap.Data[i] != -1 && b.PlayersMap.Data[i] != t.MyColor {
			b.Enemies = append(b.Enemies, m.NewPoint(i, b.Width))
		}
	}

	b.ForcesMap = m.NewMap(mapWidth)
	for i := 0; i < mapSize; i++ {
		value, _ := strconv.ParseInt(t.Forces[i*3:i*3+3], 36, 32)
		if value != 0 {
			b.ForcesMap.Data[b.rotate(i)] = int(value)
		}
	}

	filteredGold := []m.Point{}
	linq.From(b.GoldList).WhereT(func(p m.Point) bool {
		return b.PlayersMap.Get(p) != t.MyColor
	}).ToSlice(&filteredGold)

	b.GoldList = filteredGold

	// fmt.Println("Walk map")
	// b.WalkMap.Print()
	// fmt.Println("Players map")
	// b.PlayersMap.Print()
	// fmt.Println("Forces map")
	// b.ForcesMap.Print()

}

func (b *Board) buildOutsideMap() {

	points := make(map[m.Point]bool)

	b.PlayersMap.Iterate(func(i, v int) {
		if v == b.TurnInfo.MyColor {
			points[m.NewPoint(i, b.Width)] = true
		}
	})

	fmt.Printf("My forces count: %v\n", len(points))

	b.OutsideMap = b.WalkMap.Clone(func(v int) int {
		return v - 1
	})

	turn := 0
	for len(points) > 0 {
		turn++
		changes := make(map[m.Point]bool)
		for f := range points {
			if b.OutsideMap.Get(f) == 0 {
				b.OutsideMap.Set(f, turn)
			}

			b.Neighbours(f, func(pos int, p m.Point) bool {
				moveV := b.OutsideMap.Data[pos]

				if moveV == 0 {
					changes[p] = true
				}

				return true
			})
		}
		points = changes
	}
	/*
		b.MoveMap.Print()
	*/
}

func (b *Board) buildInsideMap() {

	b.InsideMap = b.OutsideMap.Clone(func(v int) int {
		switch v {
		case 2:
			return 1
		case 1:
			return 0
		default:
			return -1
		}
	})

	points := make(map[m.Point]bool)
	b.InsideMap.Iterate(func(i, v int) {
		if v == 1 {
			points[m.NewPoint(i, b.Width)] = true
		}
	})

	turn := 0
	for len(points) > 0 {
		turn++
		changes := make(map[m.Point]bool)
		for f := range points {
			if b.InsideMap.Get(f) == 0 {
				b.InsideMap.Set(f, turn)
			}

			b.Neighbours(f, func(pos int, p m.Point) bool {
				moveV := b.InsideMap.Data[pos]

				if moveV == 0 {
					changes[p] = true
				}

				return true
			})
		}
		points = changes
	}
	//b.InsideMap.Print()
}

func (b *Board) GetDistance(p m.Point) int {
	return b.WalkMap.Get(p) - 1
}

func (b *Board) GetDirection(p m.Point, pmap *m.Map) *m.Movement {
	pos := pmap.Get(p)

	dir := ""
	b.Neighbours(p, func(n_pos int, neighbour m.Point) bool {
		//found
		if pmap.Data[n_pos] < pos && pmap.Data[n_pos] > 0 {
			dir = p.GetDirection(neighbour)
			return false
		}
		return true
	})

	return &m.Movement{
		Direction: dir,
		Region:    p,
	}
}

func (b *Board) GetDirectionTo(p m.Point, pmap *m.Map) []m.Movement {
	pos := pmap.Get(p)

	if pos == 0 {
		return []m.Movement{}
	}

	moves := []m.Movement{}
	for len(moves) == 0 {
		b.Neighbours(p, func(n_pos int, neighbour m.Point) bool {
			//found
			if pmap.Data[n_pos] == 1 {
				moves = append(moves, m.Movement{
					Direction: neighbour.GetDirection(p),
					Region:    neighbour,
				})
			} else if pmap.Data[n_pos] == pos-1 {
				p = neighbour
				pos--
				return false
			}
			return true
		})
	}

	return moves
}

func (b *Board) GetDirectionFromTo(pa m.Point, pb m.Point) *m.Movement {

	pathMap := b.WalkMap.Clone(func(v int) int {
		return v - 1
	})

	points := make(map[m.Point]bool)
	points[pb] = true

	found := false

	turn := 0
	for !found {
		turn++
		changes := make(map[m.Point]bool)
		for f := range points {
			if pathMap.Get(f) == 0 {
				pathMap.Set(f, turn)
			}

			b.Neighbours(f, func(pos int, p m.Point) bool {
				moveV := pathMap.Data[pos]

				if moveV == 0 {
					changes[p] = true
				}

				if p.X == pa.X && p.Y == pa.Y {
					pathMap.Set(p, turn+1)
					found = true
					return false
				}

				return true
			})
		}
		points = changes
	}
	move := b.GetDirection(pa, pathMap)
	return move
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
