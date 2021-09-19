package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gt "github.com/buger/goterm"
	"github.com/klauspost/compress/s2"
)

// GenBuildFile creates a build file from bc.If the destination
// file already exists, it will be overwritten. If the destination
// is omitted, the build file will be saved in the CWD and will have the
// pattern's filename without the extension.
func GenBuildFile(bc BuildConfig) error {
	if bc.BuildFilePath == "" {
		var (
			patFile        = filepath.Base(bc.InitPattern.FilePath)
			fileExtLen     = len(filepath.Ext(patFile))
			trimmedPatFile = patFile[:len(patFile)-fileExtLen]
		)

		bc.BuildFilePath = "./" + trimmedPatFile
	}

	var (
		frames     = genFramesFromPattern(bc)
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

	fmt.Println("Compressing JSON...")

	framesCompressed := s2.EncodeBest(nil, framesJSON)

	fmt.Println("Writing compressed data to file...")

	err = os.WriteFile(bc.BuildFilePath, framesCompressed, 0666)
	if err != nil {
		return err
	}

	fmt.Println("Done. Saved to", bc.BuildFilePath)

	return nil
}

// RunBuildFile runs a build file according to the parameters passed in rc.
func RunBuildFile(rc RunConfig) error {
	gameData, err := genGameDataFromBuildFile(rc)
	if err != nil {
		return err
	}

	runFrames(gameData, rc)

	return nil
}

// BruteRunGame runs the game with in-time building.
func BruteRunGame(gc GameConfig) {
	var (
		bc = BuildConfig{
			Width:       gc.Width,
			Height:      gc.Height,
			FrameCount:  gc.FrameCount,
			InitPattern: gc.InitPattern,
		}
		iValComparer                          = gc.FrameCount
		frameCells                            = genCellsFromPattern(bc)
		gameConfigString, patternConfigString = getConfigListStrings(gc)
	)

	if gc.FrameCount == -1 {
		iValComparer = 9223372036854775807
	}

	gt.Clear()
	gt.MoveCursor(1, 1)

	gt.Println(gt.Color(gameConfigString, gt.YELLOW))
	gt.Println(gt.Color(patternConfigString, gt.CYAN))

	for i := 0; i < iValComparer; i++ {
		clearAndDrawFrames(gc, frameCells)

		updateCells(bc, &frameCells)

		time.Sleep(gc.Interval)
	}
}

// runFrames runs the frames based on the passed game data and run config.
func runFrames(gd GameData, rc RunConfig) {
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
		gameConfigString, patternConfigString = getConfigListStrings(gc)
	)

	gt.Clear()
	gt.MoveCursor(1, 1)

	gt.Println(gt.Color(gameConfigString, gt.YELLOW))
	gt.Println(gt.Color(patternConfigString, gt.CYAN))

	for _, frame := range gd.Frames {
		clearAndDrawFrames(gc, frame.Cells)

		time.Sleep(rc.Interval)
	}
}

// updateCells updates the value of the pointer frameCells after evaluating the living states of its cells.
func updateCells(bc BuildConfig, frameCells *Cells) {
	var (
		newFrameCells Cells
		newCell       Cell
	)

	for _, cell := range *frameCells {
		newCell = getNewCell(bc, *frameCells, cell)

		newFrameCells = append(newFrameCells, newCell)
	}

	for i := range newFrameCells {
		newFrameCells[i].LiveNeighborNum = getLiveNeighborNumByCoord(bc, newFrameCells, newFrameCells[i].X, newFrameCells[i].Y)
	}

	*frameCells = newFrameCells
}

// clearAndDrawFrames updates the frame that is printed to the command line.
func clearAndDrawFrames(gc GameConfig, frameCells Cells) {
	var (
		cellNum      int
		outputString string
	)

	for row := 0; row < gc.Height; row++ {
		for col := 0; col < gc.Width; col++ {
			cell := frameCells[cellNum]

			if cell.IsAlive {
				outputString += gc.LiveCellChar
			} else {
				outputString += gc.DeadCellChar
			}

			cellNum++
		}
		outputString += "\n"
	}

	gt.MoveCursor(1, 3)
	gt.Print(outputString)
	gt.Flush()
}
