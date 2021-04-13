package main

import (
	"fmt"
	"log"

	"github.com/gosuri/uilive"
)

func ClearAndSpawnCells(writer *uilive.Writer, frameConfig FrameConfig, frameCells FrameCells) {
	var (
		cellNum      int
		outputString string
	)

	for row := 0; row < frameConfig.Height; row++ {
		for col := 0; col < frameConfig.Width; col++ {
			cell := frameCells[cellNum]

			if cell.IsAlive {
				outputString += "#" //█
			} else {
				outputString += "-"
			}

			cellNum++
		}
		outputString += "\n"
	}

	fmt.Fprint(writer, outputString)
}

func GenerateFrameCellsFromLivingCells(frameConfig FrameConfig, initLivingCells InitLivingCells) FrameCells {
	var frameCells FrameCells

	for x := 0; x < frameConfig.Width; x++ {
		for y := 0; y < frameConfig.Height; y++ {
			cell := Cell{x, y, false}

			frameCells = append(frameCells, cell)
		}
	}

	for _, coord := range initLivingCells {
		targetCell := frameCells.GetCellByCoord(frameConfig, coord[0], coord[1])

		targetCell.IsAlive = true
	}

	return frameCells
}

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

// These rules, which compare the behavior of the automaton to real life, can be condensed into the following:

func UpdateCells(frameConfig FrameConfig, frameCells *FrameCells) {
	for _, cell := range *frameCells {
		if frameConfig.IsCoordOutOfFrame(cell.X, cell.Y) {
			log.Fatal("SetCellIsAlive: coord is out of frame")
		}

		var (
			cellNeighbors   = frameCells.GetNeighborsByCoord(frameConfig, cell.X, cell.Y)
			livingNeighbors CellNeighbors
		)

		for _, neighborCell := range cellNeighbors {
			if neighborCell.IsAlive {
				livingNeighbors = append(livingNeighbors, neighborCell)
			}
		}

		// Any live cell with two or three live neighbours survives.
		// Any dead cell with three live neighbours becomes a live cell.
		// All other live cells die in the next generation. Similarly, all other dead cells stay dead.

		// i know this part can be more concise but this is for readability
		if cell.IsAlive && (len(livingNeighbors) == 2 || len(livingNeighbors) == 3) {
			cell.IsAlive = true
		} else if !cell.IsAlive && len(livingNeighbors) == 3 {
			cell.IsAlive = true
		} else {
			cell.IsAlive = false
		}
	}
}

func (frameCells FrameCells) GetCellByCoord(frameConfig FrameConfig, x, y int) Cell {
	var targetCell Cell

	if frameConfig.IsCoordOutOfFrame(x, y) {
		log.Fatal("GetCellByCoord: coord is out of frame")
	}

	for _, cell := range frameCells {
		if cell.X == x && cell.Y == y {
			targetCell = cell
			break
		}
	}

	return targetCell
}

func (frameCells FrameCells) GetNeighborsByCoord(frameConfig FrameConfig, x, y int) CellNeighbors {
	if frameConfig.IsCoordOutOfFrame(x, y) {
		log.Fatal("GetNeighborsByCoord: coord is out of frame")
	}

	var cellNeighbors CellNeighbors

	for relX := -1; relX < 1; relX++ {
		for relY := -1; relY < 1; relY++ {
			var (
				targetX = x + relX
				targetY = y + relY
			)

			if !(relX == 0 && relY == 0) && !frameConfig.IsCoordOutOfFrame(targetX, targetY) {
				cellNeighbors = append(cellNeighbors, frameCells.GetCellByCoord(frameConfig, targetX, targetY))
			}
		}
	}

	return cellNeighbors
}

func (frameConfig FrameConfig) IsCoordOutOfFrame(x, y int) bool {
	return x < 0 || y < 0 || x > frameConfig.Width || y > frameConfig.Height
}
