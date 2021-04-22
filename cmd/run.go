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

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Runs a build golci file.",
	Args:    cobra.MaximumNArgs(1),
	Example: `golcli run ./buildFile.json`,
	Run: func(cmd *cobra.Command, args []string) {
		buildFilePath = args[0]

		frameInterval, err = cmd.Flags().GetInt("interval")
		livingCellChar, err = cmd.Flags().GetString("live-char")
		deadCellChar, err = cmd.Flags().GetString("dead-char")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gameConfig := u.GameConfig{
			Interval:       time.Duration(frameInterval) * time.Millisecond,
			LivingCellChar: livingCellChar,
			DeadCellChar:   deadCellChar,
		}

		err = u.RunBuildFile(gameConfig, buildFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&frameInterval, "interval", "i", 50, "The number of milliseconds between frames")
	runCmd.Flags().StringVarP(&livingCellChar, "live-char", "l", "██", "The character(s) that represent a live cell")
	runCmd.Flags().StringVarP(&deadCellChar, "dead-char", "d", "  ", "The character(s) that represent a live cell")
}
