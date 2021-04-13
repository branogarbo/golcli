package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gosuri/uilive"
)

func RunGame(frameConfig FrameConfig, initPattern Pattern) {
	var (
		writer     *uilive.Writer
		frameCells = GetFrameCellsByPattern(frameConfig, initPattern)
	)

	writer = uilive.New()
	writer.Start()

	for i := 0; i < frameConfig.FrameCount; i++ {
		ClearAndSpawnCells(writer, frameConfig, frameCells)

		frameCells = UpdateCells(frameConfig, frameCells)

		time.Sleep(frameConfig.Interval)
	}
}

func GetFrameCellsByPattern(frameConfig FrameConfig, pattern Pattern) FrameCells {
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

	for y := 0; y < frameConfig.Height; y++ {
		for x := 0; x < frameConfig.Width; x++ {

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
		frameCells[i].LivingNeighbors = GetLivingNeighborsByCoord(frameConfig, frameCells, frameCells[i].X, frameCells[i].Y)
	}

	return frameCells
}

func ClearAndSpawnCells(writer *uilive.Writer, frameConfig FrameConfig, frameCells FrameCells) {
	var (
		cellNum      int
		outputString string
	)

	for row := 0; row < frameConfig.Height; row++ {
		for col := 0; col < frameConfig.Width; col++ {
			cell := frameCells[cellNum]

			if cell.IsAlive {
				outputString += frameConfig.LivingCellChar
			} else {
				outputString += frameConfig.DeadCellChar
			}

			cellNum++
		}
		outputString += "\n"
	}

	fmt.Fprint(writer, outputString)
}

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

// These rules, which compare the behavior of the automaton to real life, can be condensed into the following:

func UpdateCells(frameConfig FrameConfig, frameCells FrameCells) FrameCells {
	var newFrameCells FrameCells

	for _, cell := range frameCells {
		if IsCoordOutOfFrame(frameConfig, cell.X, cell.Y) {
			log.Fatal("SetCellIsAlive: coord is out of frame")
		}

		livingNeighbors := GetLivingNeighborsByCoord(frameConfig, frameCells, cell.X, cell.Y)

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
		newFrameCells[i].LivingNeighbors = GetLivingNeighborsByCoord(frameConfig, newFrameCells, newFrameCells[i].X, newFrameCells[i].Y)
	}

	return newFrameCells
}

func GetCellByCoord(frameConfig FrameConfig, frameCells FrameCells, x, y int) Cell {
	var targetCell Cell

	if IsCoordOutOfFrame(frameConfig, x, y) {
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

func GetLivingNeighborsByCoord(frameConfig FrameConfig, frameCells FrameCells, x, y int) CellNeighbors {
	if IsCoordOutOfFrame(frameConfig, x, y) {
		log.Fatal("GetLivingNeighborsByCoord: coord is out of frame")
	}

	var cellLivingNeighbors CellNeighbors

	for relX := -1; relX <= 1; relX++ {
		for relY := -1; relY <= 1; relY++ {
			var (
				targetX = x + relX
				targetY = y + relY
			)

			if !(relX == 0 && relY == 0) && !IsCoordOutOfFrame(frameConfig, targetX, targetY) {
				cell := GetCellByCoord(frameConfig, frameCells, targetX, targetY)

				if cell.IsAlive {
					cellLivingNeighbors = append(cellLivingNeighbors, cell)
				}
			}
		}
	}

	return cellLivingNeighbors
}

func IsCoordOutOfFrame(frameConfig FrameConfig, x, y int) bool {
	return x < 0 || y < 0 || x > frameConfig.Width || y > frameConfig.Height
}
