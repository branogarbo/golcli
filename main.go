package main

import (
	"time"
)

func main() {
	var (
		frameConfig = FrameConfig{
			Width:          70,
			Height:         50,
			FrameCount:     99999999,
			Interval:       10 * time.Millisecond,
			LivingCellChar: "██",
			DeadCellChar:   "  ",
		}
		initPattern = Pattern{
			Path: "./initPatterns/spaceshipCrash.txt",
			X:    24,
			Y:    18,
		}
	)

	RunGame(frameConfig, initPattern)
}
