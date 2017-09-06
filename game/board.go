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

	/*
		fmt.Println("Walk map")
		b.WalkMap.Print()
		fmt.Println("Players map")
		b.PlayersMap.Print()
		fmt.Println("Forces map")
		b.ForcesMap.Print()
	*/
}

func (b *Board) buildOutsideMap() {

	points := []m.Point{}

	b.PlayersMap.Iterate(func(i, v int) {
		if v == b.TurnInfo.MyColor {
			points = append(points, m.NewPoint(i, b.Width))
		}
	})

	fmt.Printf("My forces: %v\n", points)

	b.OutsideMap = b.WalkMap.Clone(func(v int) int {
		return v - 1
	})

	turn := 0
	for len(points) > 0 {
		turn++
		changes := []m.Point{}
		for _, f := range points {
			if b.OutsideMap.Get(f) == 0 {
				b.OutsideMap.Set(f, turn)
			}

			b.Neighbours(f, func(pos int, p m.Point) bool {
				moveV := b.OutsideMap.Data[pos]

				if moveV == 0 {
					changes = append(changes, p)
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

	points := []m.Point{}
	b.InsideMap.Iterate(func(i, v int) {
		if v == 1 {
			points = append(points, m.NewPoint(i, b.Width))
		}
	})

	turn := 0
	for len(points) > 0 {
		turn++
		changes := []m.Point{}
		for _, f := range points {
			if b.InsideMap.Get(f) == 0 {
				b.InsideMap.Set(f, turn)
			}

			b.Neighbours(f, func(pos int, p m.Point) bool {
				moveV := b.InsideMap.Data[pos]

				if moveV == 0 {
					changes = append(changes, p)
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

func (b *Board) GetDirectionTo(p m.Point, pmap *m.Map) *m.Movement {
	pos := pmap.Get(p)

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
			if pmap.Data[n_pos] == 1 {
				dir = neighbour.GetDirection(p)
				p = neighbour
				return false
			} else if pmap.Data[n_pos] == pos-1 {
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
