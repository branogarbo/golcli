/*
Copyright © 2021 Brian Longmore brianl.ext@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	u "github.com/branogarbo/golcli/util"
	"github.com/spf13/cobra"
)

var (
	frameWidth     int
	frameHeight    int
	frameCount     int
	frameInterval  int
	livingCellChar string
	deadCellChar   string
	patternPath    string
	patternX       int
	patternY       int
	err            error
)

var rootCmd = &cobra.Command{
	Use:   "golcli",
	Short: "A simple CLI implementation of Conway's Game of Life.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			patternPath = "./pattern.txt"
		} else {
			patternPath = args[0]
		}

		frameWidth, err = cmd.Flags().GetInt("width")
		frameHeight, err = cmd.Flags().GetInt("height")
		frameCount, err = cmd.Flags().GetInt("count")
		frameInterval, err = cmd.Flags().GetInt("interval")
		livingCellChar, err = cmd.Flags().GetString("live-char")
		deadCellChar, err = cmd.Flags().GetString("dead-char")

		patternX, err = cmd.Flags().GetInt("pattern-x")
		patternY, err = cmd.Flags().GetInt("pattern-y")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gameConfig := u.GameConfig{
			Width:          frameWidth,
			Height:         frameHeight,
			FrameCount:     frameCount,
			Interval:       time.Duration(frameInterval) * time.Millisecond,
			LivingCellChar: livingCellChar,
			DeadCellChar:   deadCellChar,
		}

		initPattern := u.Pattern{
			Path: patternPath,
			X:    patternX,
			Y:    patternY,
		}

		u.RunGame(gameConfig, initPattern)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&frameWidth, "width", "W", 40, "The width of the frames.")
	rootCmd.Flags().IntVarP(&frameHeight, "height", "H", 30, "The height of the frames.")
	rootCmd.Flags().IntVarP(&frameCount, "count", "c", 999999999, "The number of frames displayed before exiting.")
	rootCmd.Flags().IntVarP(&frameInterval, "interval", "i", 50, "The number of milliseconds between frames.")
	rootCmd.Flags().StringVarP(&livingCellChar, "live-char", "l", "██", "The character(s) that represent a live cell.")
	rootCmd.Flags().StringVarP(&deadCellChar, "dead-char", "d", "  ", "The character(s) that represent a live cell.")

	rootCmd.Flags().IntVarP(&patternX, "pattern-x", "x", 12, "The initial x offset of the pattern defined with --pattern-file.")
	rootCmd.Flags().IntVarP(&patternY, "pattern-y", "y", 8, "The initial y offset of the pattern defined with --pattern-file.")
}
