package models

type TurnInfo struct {
	OnlyMyName    bool
	MyColor       int
	ShowName      bool
	Offset        Point
	MyBase        Point
	Forces        string
	Layers        []string
	LevelProgress LevelProgress
}
