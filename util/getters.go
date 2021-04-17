package util

import (
	"log"
	"os"
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

		strX++

		if char == "#" {
			initLivingCells = append(initLivingCells, [2]int{strX + pattern.X, strY + pattern.Y})
		} else if char == "\n" {
			strX = 0
			strY++
		}

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
		frameCells[i].LivingNeighbors = GetLivingNeighborsByCoord(gameConfig, frameCells, frameCells[i].X, frameCells[i].Y)
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
func GetLivingNeighborsByCoord(gameConfig GameConfig, frameCells FrameCells, x, y int) CellNeighbors {
	if IsCoordOutOfFrame(gameConfig, x, y) {
		log.Fatal("GetLivingNeighborsByCoord: coord is out of frame")
	}

	var cellLivingNeighbors CellNeighbors

	for relX := -1; relX <= 1; relX++ {
		for relY := -1; relY <= 1; relY++ {
			var (
				targetX = x + relX
				targetY = y + relY
			)

			if !(relX == 0 && relY == 0) && !IsCoordOutOfFrame(gameConfig, targetX, targetY) {
				cell := GetCellByCoord(gameConfig, frameCells, targetX, targetY)

				if cell.IsAlive {
					cellLivingNeighbors = append(cellLivingNeighbors, cell)
				}
			}
		}
	}

	return cellLivingNeighbors
}

// IsCoordOutOfFrame returns whether or not the cell at (x,y) is out of the frame.
func IsCoordOutOfFrame(gameConfig GameConfig, x, y int) bool {
	return x < 0 || y < 0 || x > gameConfig.Width || y > gameConfig.Height
}
