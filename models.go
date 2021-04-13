package main

type FrameCells []Cell

type CellNeighbors []Cell

type InitLivingCells [][2]int

type FrameConfig struct {
	Width  int
	Height int
	Cells  FrameCells
}

type Cell struct {
	X               int
	Y               int
	IsAlive         bool
	LivingNeighbors CellNeighbors
}
