package models

type PlayerInfo struct {
	ActiveForces  int
	TerritorySize int
	TotalForces   int
	MinesCount    int

	Mines  map[Point]bool
	Forces map[Point]int
}
