package util

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	gt "github.com/buger/goterm"
)

// GenBuildFile creates a build file from bc.If the destination
// file already exists, it will be overwritten.
func (bc BuildConfig) GenBuildFile() error {
	var (
		frames     = bc.GenFramesFromPattern()
		framesJSON []byte
		err        error
	)

	fmt.Println("Encoding to JSON...")

	framesJSON, err = json.Marshal(GameData{
		BuildConfig: bc,
		Frames:      frames,
	})
	if err != nil {
		return err
	}

	fmt.Println("Writing JSON to file...")

	err = os.WriteFile(bc.BuildFilePath, framesJSON, 0666)
	if err != nil {
		return err
	}

	fmt.Println("Done. Saved to", bc.BuildFilePath)

	return nil
}

// RunBuildFile runs a build file according to the parameters passed in rc.
func (rc RunConfig) RunBuildFile() error {
	fmt.Println("Loading Game Data...")

	gameData, err := rc.GenGameDataFromBuildFile()
	if err != nil {
		return err
	}

	rc.RunFrames(gameData)

	return nil
}

// RunFrames runs the frames based on the passed game data.
func (rc RunConfig) RunFrames(gd GameData) {
	var (
		bc = gd.BuildConfig
		gc = GameConfig{
			Width:        bc.Width,
			Height:       bc.Height,
			FrameCount:   bc.FrameCount,
			InitPattern:  bc.InitPattern,
			Interval:     rc.Interval,
			LiveCellChar: rc.LiveCellChar,
			DeadCellChar: rc.DeadCellChar,
		}
		gameConfigString, patternConfigString = gc.GetConfigListStrings()
	)

	gt.Clear()
	gt.MoveCursor(1, 1)

	gt.Println(gt.Color(gameConfigString, gt.YELLOW))
	gt.Println(gt.Color(patternConfigString, gt.CYAN))

	for _, frame := range gd.Frames {
		rc.ClearAndDrawFrames(gd, frame.Cells)

		time.Sleep(rc.Interval)
	}
}

// UpdateCells updates the value of the pointer frameCells after evaluating the living states of its cells.
func UpdateCells(bc BuildConfig, frameCells *Cells) {
	var (
		newFrameCells Cells
		newCell       Cell
	)

	for _, cell := range *frameCells {
		newCell = bc.GetNewCell(*frameCells, cell)

		newFrameCells = append(newFrameCells, newCell)
	}

	for i := range newFrameCells {
		newFrameCells[i].LiveNeighborNum = bc.GetLiveNeighborNumByCoord(newFrameCells, newFrameCells[i].X, newFrameCells[i].Y)
	}

	*frameCells = newFrameCells
}

// ClearAndDrawFrames updates the frame that is printed to the command line.
func (rc RunConfig) ClearAndDrawFrames(gd GameData, frameCells Cells) {
	var (
		cellNum      int
		outputString string
		bc           = gd.BuildConfig
	)

	for row := 0; row < bc.Height; row++ {
		for col := 0; col < bc.Width; col++ {
			cell := frameCells[cellNum]

			if cell.IsAlive {
				outputString += rc.LiveCellChar
			} else {
				outputString += rc.DeadCellChar
			}

			cellNum++
		}
		outputString += "\n"
	}

	gt.MoveCursor(1, 3)
	gt.Print(outputString)
	gt.Flush()
}
