package util

import (
	"fmt"
	"log"
	"time"

	"github.com/gosuri/uilive"
)

// RunGame runs the game based on the passed args.
func RunGame(gameConfig GameConfig, initPattern Pattern) {
	var (
		writer       *uilive.Writer
		iValComparer = gameConfig.FrameCount
		frameCells   = GetFrameCellsByPattern(gameConfig, initPattern)
	)

	writer = uilive.New()
	writer.Start()

	if gameConfig.FrameCount == -1 {
		iValComparer = 9223372036854775807
	}

	for i := 0; i < iValComparer; i++ {
		ClearAndSpawnCells(writer, gameConfig, frameCells)

		frameCells = UpdateCells(gameConfig, frameCells)

		time.Sleep(gameConfig.Interval)
	}
}

// ClearAndSpawnCells updates the cell string that is printed to the command line.
func ClearAndSpawnCells(writer *uilive.Writer, gameConfig GameConfig, frameCells FrameCells) {
	var (
		cellNum      int
		outputString string
	)

	for row := 0; row < gameConfig.Height; row++ {
		for col := 0; col < gameConfig.Width; col++ {
			cell := frameCells[cellNum]

			if cell.IsAlive {
				outputString += gameConfig.LivingCellChar
			} else {
				outputString += gameConfig.DeadCellChar
			}

			cellNum++
		}
		outputString += "\n"
	}

	fmt.Fprint(writer, outputString)
}

// UpdateCells returns new frame cells after evaluating the living state of frameCells.
func UpdateCells(gameConfig GameConfig, frameCells FrameCells) FrameCells {
	var newFrameCells FrameCells

	for _, cell := range frameCells {
		if IsCoordOutOfFrame(gameConfig, cell.X, cell.Y) {
			log.Fatal("SetCellIsAlive: coord is out of frame")
		}

		livingNeighbors := GetLivingNeighborsByCoord(gameConfig, frameCells, cell.X, cell.Y)

		// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
		// Any live cell with two or three live neighbours lives on to the next generation.
		// Any live cell with more than three live neighbours dies, as if by overpopulation.
		// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

		// These rules, which compare the behavior of the automaton to real life, can be condensed into the following:

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

		newFrameCells = append(newFrameCells, cell)
	}

	for i := range newFrameCells {
		newFrameCells[i].LivingNeighbors = GetLivingNeighborsByCoord(gameConfig, newFrameCells, newFrameCells[i].X, newFrameCells[i].Y)
	}

	return newFrameCells
}
