package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/klauspost/compress/s2"
)

// GenCellsFromPattern converts a pattern string to frame cells.
func GenCellsFromPattern(bc BuildConfig) Cells {
	var (
		frameCells      Cells
		initLivingCells [][2]int
		isAlive         bool
		strX            int
		strY            int
		pattern         = bc.InitPattern
	)

	fileBytes, err := os.ReadFile(pattern.FilePath)
	if err != nil {
		log.Fatal(err)
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

	for y := 0; y < bc.Height; y++ {
		for x := 0; x < bc.Width; x++ {

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
		frameCells[i].LiveNeighborNum = GetLiveNeighborNumByCoord(bc, frameCells, frameCells[i].X, frameCells[i].Y)
	}

	return frameCells
}

// GetCellByCoord returns the cell that is at (x,y).
func GetCellByCoord(bc BuildConfig, frameCells Cells, x, y int) Cell {
	var targetCell Cell

	if IsCoordOutOfFrame(bc, x, y) {
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

// GetLiveNeighborNumByCoord returns the number of cells neighboring the cell at (x,y).
func GetLiveNeighborNumByCoord(bc BuildConfig, frameCells Cells, x, y int) int {
	if IsCoordOutOfFrame(bc, x, y) {
		log.Fatal("GetLivingNeighborsByCoord: coord is out of frame")
	}

	var liveNeighborNum int

	for relX := -1; relX <= 1; relX++ {
		for relY := -1; relY <= 1; relY++ {
			var (
				targetX = x + relX
				targetY = y + relY
			)

			if !(relX == 0 && relY == 0) && !IsCoordOutOfFrame(bc, targetX, targetY) {
				cell := GetCellByCoord(bc, frameCells, targetX, targetY)

				if cell.IsAlive {
					liveNeighborNum++
				}
			}
		}
	}

	return liveNeighborNum
}

// IsCoordOutOfFrame returns whether or not the cell at (x,y) is out of the frame.
func IsCoordOutOfFrame(bc BuildConfig, x, y int) bool {
	return x < 0 || y < 0 || x > bc.Width || y > bc.Height
}

// GetNewCell returns a new cell from the passed cell according to the life conditions.
func GetNewCell(bc BuildConfig, frameCells Cells, cell Cell) Cell {
	if IsCoordOutOfFrame(bc, cell.X, cell.Y) {
		log.Fatal("GetNewCell: coord is out of frame")
	}

	livingNeighbors := GetLiveNeighborNumByCoord(bc, frameCells, cell.X, cell.Y)

	// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
	// Any live cell with two or three live neighbours lives on to the next generation.
	// Any live cell with more than three live neighbours dies, as if by overpopulation.
	// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

	// These rules, which compare the behavior of the automaton to real life, can be condensed into the following:

	// Any live cell with two or three live neighbours survives.
	// Any dead cell with three live neighbours becomes a live cell.
	// All other live cells die in the next generation. Similarly, all other dead cells stay dead.

	// Rules compacted to condition
	cell.IsAlive = (cell.IsAlive && (livingNeighbors == 2 || livingNeighbors == 3)) || (!cell.IsAlive && livingNeighbors == 3)

	return cell
}

// GetConfigListStrings returns both game and pattern config strings from config.
func GetConfigListStrings(config GameConfig) (string, string) {
	var (
		gameConfigList    string
		patternConfigList string
		clv               = reflect.ValueOf(config)
		typeOfCLV         = clv.Type()
	)

	for i := 0; i < clv.NumField(); i++ {
		var (
			key     = typeOfCLV.Field(i).Name
			value   = clv.Field(i).Interface()
			valType = reflect.TypeOf(value).Kind()
		)

		if key == "InitPattern" {
			continue
		}

		if valType == reflect.String {
			value = fmt.Sprintf(`"%v"`, value)
		}

		gameConfigList += fmt.Sprintf("%v: %v  |  ", key, value)
	}

	clv = reflect.ValueOf(config.InitPattern)
	typeOfCLV = clv.Type()

	for i := 0; i < clv.NumField(); i++ {
		var (
			key   = typeOfCLV.Field(i).Name
			value = clv.Field(i).Interface()
		)

		patternConfigList += fmt.Sprintf("%v: %v  |  ", key, value)
	}

	return gameConfigList, patternConfigList
}

// GenFramesFromPattern generates and returns Frames according to the values passed in bc.
func GenFramesFromPattern(bc BuildConfig) Frames {
	var (
		frameCells  = GenCellsFromPattern(bc)
		frames      = Frames{Frame{0, frameCells}}
		pbTemplate  = `Loading frames: {{ etime . }} {{ bar . "[" "=" ">" " " "]" }} {{speed . }} {{percent . }}`
		progressBar = pb.ProgressBarTemplate(pbTemplate).Start(bc.FrameCount).SetMaxWidth(100)
	)

	for i := 0; i < bc.FrameCount; i++ {
		UpdateCells(bc, &frameCells)
		frames = append(frames, Frame{i + 1, frameCells})

		progressBar.Increment()

	}

	progressBar.Finish()

	return frames
}

// GenGameDataFromBuildFile returns GameData that was stored in the targeted json build file.
func GenGameDataFromBuildFile(rc RunConfig) (GameData, error) {
	fmt.Println("Reading build file...")

	var (
		gameData            GameData
		compressedJSON, err = os.ReadFile(rc.BuildFilePath)
	)
	if err != nil {
		return GameData{}, err
	}

	fmt.Println("Extracting JSON...")

	jsonBytes, err := s2.Decode(nil, compressedJSON)
	if err != nil {
		return GameData{}, err
	}

	err = json.Unmarshal(jsonBytes, &gameData)
	if err != nil {
		return GameData{}, err
	}

	return gameData, nil
}
