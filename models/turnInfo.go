package models

type TurnInfo struct {
	Tick          int
	Available     int
	OnlyMyName    bool
	MyColor       int
	ShowName      bool
	Offset        Point
	MyBase        Point
	Forces        string
	Layers        []string
	LevelProgress LevelProgress
}
