package main

import (
	"github.com/gosuri/uilive"
)

func main() {
	var (
		writer      *uilive.Writer
		frameConfig = FrameConfig{
			Width:  100,
			Height: 40,
		}
		initLivingCells = InitLivingCells{{20, 20}, {21, 21}, {20, 21}}
		frameCells      = GenerateFrameCellsFromLivingCells(frameConfig, initLivingCells)
	)

	writer = uilive.New()
	writer.Start()

	ClearAndSpawnCells(writer, frameConfig, frameCells)
	// for i := 0; i < 5; i++ {

	// 	frameCells = UpdateCells(frameConfig, frameCells)

	// 	time.Sleep(500 * time.Millisecond)
	// }

	writer.Stop()
}
