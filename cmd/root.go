/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	Short: "A simple implementation of Conway's Game of Life ",
	Run: func(cmd *cobra.Command, args []string) {

		frameWidth, err = cmd.Flags().GetInt("width")
		frameHeight, err = cmd.Flags().GetInt("height")
		frameCount, err = cmd.Flags().GetInt("count")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var (
			frameConfig = u.FrameConfig{
				Width:          frameWidth,
				Height:         frameHeight,
				FrameCount:     frameCount,
				Interval:       time.Duration(frameInterval) * time.Millisecond,
				LivingCellChar: livingCellChar,
				DeadCellChar:   deadCellChar,
			}
			initPattern = u.Pattern{
				Path: patternPath,
				X:    patternX,
				Y:    patternY,
			}
		)

		u.RunGame(frameConfig, initPattern)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&frameWidth, "width", "w", 70, "The width of the viewing frame.")
	rootCmd.Flags().IntVarP(&frameHeight, "height", "h", 50, "The height of the viewing frame.")
	rootCmd.Flags().IntVarP(&frameCount, "count", "c", 999999999, "The height of the viewing frame.")
	rootCmd.Flags().IntVarP(&frameInterval, "interval", "i", 10, "The number of milliseconds between frames.")
	rootCmd.Flags().IntVarP(&frameInterval, "interval", "i", 10, "The number of milliseconds between frames.")
	rootCmd.Flags().StringVarP(&livingCellChar, "live-char", "lc", "██", "The character(s) that represent a live cell.")
	rootCmd.Flags().StringVarP(&livingCellChar, "dead-char", "dc", "  ", "The character(s) that represent a live cell.")

	rootCmd.Flags().StringVarP(&patternPath, "pattern-file", "-p", "./pattern.txt", `The initial pattern of live and dead cells, use "#" to represent live cells.`)
	rootCmd.Flags().IntVarP(&patternX, "pattern-x", "-px", 24, `The initial x offset of the pattern defined with --pattern-file.`)
	rootCmd.Flags().IntVarP(&patternY, "pattern-y", "-py", 18, `The initial y offset of the pattern defined with --pattern-file.`)
}
