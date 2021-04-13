package main

import "time"

type FrameCells []Cell

type CellNeighbors []Cell

type Pattern struct {
	Path string
	X    int
	Y    int
}

type FrameConfig struct {
	Width          int
	Height         int
	FrameCount     int
	Interval       time.Duration
	DeadCellChar   string
	LivingCellChar string
}

type Cell struct {
	X               int
	Y               int
	IsAlive         bool
	LivingNeighbors CellNeighbors
}
