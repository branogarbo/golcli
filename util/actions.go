package util

import (
	"fmt"
	"reflect"
	"time"

	gt "github.com/buger/goterm"
)

// RunGame runs the game based on the passed args.
func RunGame(gameConfig GameConfig, initPattern Pattern) {
	var (
		iValComparer = gameConfig.FrameCount
		frameCells   = GetFrameCellsByPattern(gameConfig, initPattern)
	)

	if gameConfig.FrameCount == -1 {
		iValComparer = 9223372036854775807
	}

	gt.Clear()

	for i := 0; i < iValComparer; i++ {
		ClearAndSpawnCells(gameConfig, frameCells)

		frameCells = UpdateCells(gameConfig, frameCells)

		time.Sleep(gameConfig.Interval)
	}
}

func GetConfigListString(gameConfig GameConfig) string {
	var (
		configList string
		v          = reflect.ValueOf(gameConfig)
		typeOfGC   = v.Type()
	)

	for i := 0; i < v.NumField(); i++ {
		var (
			key     = typeOfGC.Field(i).Name
			value   = v.Field(i).Interface()
			valType = reflect.TypeOf(value).Kind()
		)

		if valType == reflect.String {
			value = fmt.Sprintf(`"%v"`, value)
		}

		configList += fmt.Sprintf("%v: %v  |  ", key, value)
	}

	return configList
}

// ClearAndSpawnCells updates the cell string that is printed to the command line.
func ClearAndSpawnCells(gameConfig GameConfig, frameCells FrameCells) {
	var (
		cellNum      int
		outputString string
		configList   string = GetConfigListString(gameConfig)
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

	gt.MoveCursor(1, 1)

	gt.Println(gt.Color(configList, gt.YELLOW))
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
