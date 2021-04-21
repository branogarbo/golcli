package util

import (
	"time"

	gt "github.com/buger/goterm"
)

// RunGame runs the game based on the passed args.
func RunGame(gameConfig GameConfig, initPattern Pattern) {
	var (
		iValComparer                          = gameConfig.FrameCount
		frameCells                            = GetFrameCellsByPattern(gameConfig, initPattern)
		gameConfigString, patternConfigString = GetConfigListStrings(gameConfig, initPattern)
	)

	if gameConfig.FrameCount == -1 {
		iValComparer = 9223372036854775807
	}

	gt.Clear()

	gt.MoveCursor(1, 1)

	gt.Println(gt.Color(gameConfigString, gt.YELLOW))
	gt.Println(gt.Color(patternConfigString, gt.CYAN))

	for i := 0; i < iValComparer; i++ {
		ClearAndSpawnCells(gameConfig, frameCells)

		UpdateCells(gameConfig, &frameCells)

		time.Sleep(gameConfig.Interval)
	}
}

// ClearAndSpawnCells updates the cell string that is printed to the command line.
func ClearAndSpawnCells(gameConfig GameConfig, frameCells FrameCells) {
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

	gt.MoveCursor(1, 3)

	gt.Print(outputString)

	gt.Flush()
}

// UpdateCells updates the value of the pointer frameCells after evaluating the living states of its cells.
func UpdateCells(gameConfig GameConfig, frameCells *FrameCells) {
	var (
		newFrameCells FrameCells
		newCell       Cell
	)

	for _, cell := range *frameCells {
		newCell = GetNewCell(gameConfig, *frameCells, cell)

		newFrameCells = append(newFrameCells, newCell)
	}

	for i := range newFrameCells {
		newFrameCells[i].LivingNeighbors = GetLivingNeighborsByCoord(gameConfig, newFrameCells, newFrameCells[i].X, newFrameCells[i].Y)
	}

	*frameCells = newFrameCells
}
