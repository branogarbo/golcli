package util

import (
	"fmt"
	"log"
	"os"
	"reflect"

	pb "github.com/cheggaaa/pb/v3"
)

// GetFrameCellsByPattern converts a pattern string to frame cells.
func GetFrameCellsByPattern(gameConfig GameConfig, pattern Pattern) FrameCells {
	var (
		frameCells      FrameCells
		initLivingCells [][2]int
		isAlive         bool
		strX            int
		strY            int
	)

	fileBytes, err := os.ReadFile(pattern.Path)
	if err != nil {
		log.Fatal("ConvPatternToFrameCells: could not open pattern file")
	}

	for _, c := range fileBytes {
		char := string(c)

		if char == "#" {
			initLivingCells = append(initLivingCells, [2]int{strX + pattern.X, strY + pattern.Y})
		} else if char == "\n" {
			strX = -1
			strY++
		}

		strX++
	}

	for y := 0; y < gameConfig.Height; y++ {
		for x := 0; x < gameConfig.Width; x++ {

			for _, coord := range initLivingCells {
				isAlive = x == coord[0] && y == coord[1]

				if isAlive {
					break
				}
			}

			cell := Cell{
				X:       x,
				Y:       y,
				IsAlive: isAlive,
			}

			frameCells = append(frameCells, cell)
		}
	}

	for i := range frameCells {
		frameCells[i].LivingNeighborsNum = GetLivingNeighborsNumByCoord(gameConfig, frameCells, frameCells[i].X, frameCells[i].Y)
	}

	return frameCells
}

// GetCellByCoord returns the cell that is at (x,y).
func GetCellByCoord(gameConfig GameConfig, frameCells FrameCells, x, y int) Cell {
	var targetCell Cell

	if IsCoordOutOfFrame(gameConfig, x, y) {
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

// GetLivingNeighborsByCoord returns a slice of cells that neighbor the cell at (x,y).
func GetLivingNeighborsNumByCoord(gameConfig GameConfig, frameCells FrameCells, x, y int) int {
	if IsCoordOutOfFrame(gameConfig, x, y) {
		log.Fatal("GetLivingNeighborsByCoord: coord is out of frame")
	}

	var livingNeighborsNum int

	for relX := -1; relX <= 1; relX++ {
		for relY := -1; relY <= 1; relY++ {
			var (
				targetX = x + relX
				targetY = y + relY
			)

			if !(relX == 0 && relY == 0) && !IsCoordOutOfFrame(gameConfig, targetX, targetY) {
				cell := GetCellByCoord(gameConfig, frameCells, targetX, targetY)

				if cell.IsAlive {
					livingNeighborsNum++
				}
			}
		}
	}

	return livingNeighborsNum
}

// IsCoordOutOfFrame returns whether or not the cell at (x,y) is out of the frame.
func IsCoordOutOfFrame(gameConfig GameConfig, x, y int) bool {
	return x < 0 || y < 0 || x > gameConfig.Width || y > gameConfig.Height
}

// GetNewCell returns a new cell from the passed cell according to the life conditions.
func GetNewCell(gameConfig GameConfig, frameCells FrameCells, cell Cell) Cell {
	if IsCoordOutOfFrame(gameConfig, cell.X, cell.Y) {
		log.Fatal("GetNewCell: coord is out of frame")
	}

	livingNeighbors := GetLivingNeighborsNumByCoord(gameConfig, frameCells, cell.X, cell.Y)

	// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
	// Any live cell with two or three live neighbours lives on to the next generation.
	// Any live cell with more than three live neighbours dies, as if by overpopulation.
	// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

	// These rules, which compare the behavior of the automaton to real life, can be condensed into the following:

	// Any live cell with two or three live neighbours survives.
	// Any dead cell with three live neighbours becomes a live cell.
	// All other live cells die in the next generation. Similarly, all other dead cells stay dead.

	// i know this part can be more concise but this is for readability
	if cell.IsAlive && (livingNeighbors == 2 || livingNeighbors == 3) {
		cell.IsAlive = true
	} else if !cell.IsAlive && livingNeighbors == 3 {
		cell.IsAlive = true
	} else {
		cell.IsAlive = false
	}

	return cell
}

// GetConfigListString returns config list string from config.
func GetConfigListString(config interface{}) string {
	var (
		configList string
		clv        = reflect.ValueOf(config)
		typeOfCLV  = clv.Type()
	)

	for i := 0; i < clv.NumField(); i++ {
		var (
			key     = typeOfCLV.Field(i).Name
			value   = clv.Field(i).Interface()
			valType = reflect.TypeOf(value).Kind()
		)

		if valType == reflect.String {
			value = fmt.Sprintf(`"%v"`, value)
		}

		configList += fmt.Sprintf("%v: %v  |  ", key, value)
	}

	return configList
}

// GenerateFrames returns a chan of frameCells that represent each frame in the game.
func GenFramesFromPattern(gameConfig GameConfig, initPattern Pattern) Frames {
	var (
		frameCells  = GetFrameCellsByPattern(gameConfig, initPattern)
		frames      = Frames{frameCells}
		pbTemplate  = `Loading frames: {{ etime . }} {{ bar . "[" "=" ">" " " "]" }} {{speed . }} {{percent . }}`
		progressBar = pb.ProgressBarTemplate(pbTemplate).Start(gameConfig.FrameCount).SetMaxWidth(100)
	)

	for i := 0; i < gameConfig.FrameCount; i++ {
		UpdateCells(gameConfig, &frameCells)
		frames = append(frames, frameCells)

		progressBar.Increment()

	}

	progressBar.Finish()

	return frames
}
