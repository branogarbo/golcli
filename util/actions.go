package util

import (
	"time"

	gt "github.com/buger/goterm"
	pb "github.com/cheggaaa/pb/v3"
)

// RunGame runs the game based on the passed args.
func RunGame(gameConfig GameConfig, initPattern Pattern) {
	var (
		frameCells                            = GetFrameCellsByPattern(gameConfig, initPattern)
		gameConfigString, patternConfigString = GetConfigListStrings(gameConfig, initPattern)
		pbTemplate                            = `{{ etime . }} {{ bar . "[" "=" ">" " " "]" }} {{speed . }} {{percent . }}`
		progressBar                           = pb.ProgressBarTemplate(pbTemplate).Start(gameConfig.FrameCount).SetMaxWidth(100)
	)

	for i := 0; i < gameConfig.FrameCount; i++ {
		// go func() {
		frameCells = UpdateCells(gameConfig, frameCells)
		progressBar.Increment()
		// }()
	}

	progressBar.Finish()

	gt.Clear()

	gt.MoveCursor(1, 1)

	gt.Println(gt.Color(gameConfigString, gt.YELLOW))
	gt.Println(gt.Color(patternConfigString, gt.CYAN))

	for i := 0; i < gameConfig.FrameCount; i++ {
		ClearAndSpawnCells(gameConfig, frameCells)

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

// UpdateCells returns new frame cells after evaluating the living state of frameCells.
func UpdateCells(gameConfig GameConfig, frameCells FrameCells) FrameCells {
	var (
		newFrameCells FrameCells
		newCell       Cell
	)

	for _, cell := range frameCells {
		newCell = GetNewCell(gameConfig, frameCells, cell)

		newFrameCells = append(newFrameCells, newCell)
	}

	for i := range newFrameCells {
		newFrameCells[i].LivingNeighbors = GetLivingNeighborsByCoord(gameConfig, newFrameCells, newFrameCells[i].X, newFrameCells[i].Y)
	}

	return newFrameCells
}
