/*
Copyright © 2021 Brian Longmore branodev@gmail.com

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
	"time"

	u "github.com/branogarbo/golcli/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "golcli",
	Short:   "A basic CLI implementation of Conway's Game of Life.",
	Example: "golcli -c 100 -i 20 ./pattern.txt",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		patternFilePath := args[0]

		frameWidth, _ := cmd.Flags().GetInt("width")
		frameHeight, _ := cmd.Flags().GetInt("height")
		frameCount, _ := cmd.Flags().GetInt("count")
		patternX, _ := cmd.Flags().GetInt("pattern-x")
		patternY, _ := cmd.Flags().GetInt("pattern-y")
		frameInterval, _ := cmd.Flags().GetInt("interval")
		livingCellChar, _ := cmd.Flags().GetString("live-char")
		deadCellChar, _ := cmd.Flags().GetString("dead-char")

		gc := u.GameConfig{
			Width:        frameWidth,
			Height:       frameHeight,
			FrameCount:   frameCount,
			Interval:     time.Duration(frameInterval) * time.Millisecond,
			LiveCellChar: livingCellChar,
			DeadCellChar: deadCellChar,
			InitPattern: u.Pattern{
				FilePath: patternFilePath,
				X:        patternX,
				Y:        patternY,
			},
		}

		u.BruteRunGame(gc)

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntP("interval", "i", 50, "The number of milliseconds between frames")
	rootCmd.Flags().IntP("width", "W", 40, "The width of the frames")
	rootCmd.Flags().IntP("height", "H", 30, "The height of the frames")
	rootCmd.Flags().IntP("count", "c", -1, "The number of frames displayed before exiting (-1 : infinite loop)")
	rootCmd.Flags().IntP("pattern-x", "x", 12, "The x offset of the initial pattern")
	rootCmd.Flags().IntP("pattern-y", "y", 8, "The y offset of the initial pattern")
	rootCmd.Flags().StringP("live-char", "l", "  ", "The character(s) that represent a live cell")
	rootCmd.Flags().StringP("dead-char", "d", "██", "The character(s) that represent a dead cell")
}
