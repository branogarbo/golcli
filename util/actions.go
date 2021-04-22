package util

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	gt "github.com/buger/goterm"
)

// RunGame runs the game based on the passed args.
func RunFrames(gameConfig GameConfig, frames Frames) {
	var (
		gameConfigString = GetConfigListString(gameConfig)
	)

	gt.Clear()

	gt.MoveCursor(1, 1)

	gt.Println(gt.Color(gameConfigString, gt.YELLOW))

	for _, frameCells := range frames {
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
		newFrameCells[i].LivingNeighborsNum = GetLivingNeighborsNumByCoord(gameConfig, newFrameCells, newFrameCells[i].X, newFrameCells[i].Y)
	}

	*frameCells = newFrameCells
}

// GenBuildFile creates a build file from the initPattern and gameConfig.
// If the destination file already exists, it will be overwritten.
func GenBuildFile(gameConfig GameConfig, initPattern Pattern, buildFile string) error {
	var (
		frames     = GenFramesFromPattern(gameConfig, initPattern)
		framesJSON []byte
		err        error
	)

	fmt.Println("Encoding to JSON...")

	framesJSON, err = json.Marshal(frames)
	if err != nil {
		return err
	}

	fmt.Println("Writing JSON to file...")

	err = os.WriteFile(buildFile, framesJSON, 0666)
	if err != nil {
		return err
	}

	fmt.Println("Done. Saved to", buildFile)

	return nil
}

//
func GenFramesFromBuildFile(gameConfig GameConfig, buildFile string) (Frames, error) {
	var (
		frames         Frames
		jsonBytes, err = os.ReadFile(buildFile)
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &frames)
	if err != nil {
		return nil, err
	}

	return frames, nil
}

func RunBuildFile(gameConfig GameConfig, buildFilePath string) error {
	var (
		frames, err = GenFramesFromBuildFile(gameConfig, buildFilePath)
	)
	if err != nil {
		return err
	}

	RunFrames(gameConfig, frames)

	return nil
}
