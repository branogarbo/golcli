package cmd

import (
	"fmt"
	"os"
	"time"

	u "github.com/branogarbo/golcli/util"
	"github.com/spf13/cobra"
)

var (
	frameInterval  int
	livingCellChar string
	deadCellChar   string
)

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Runs a build file.",
	Example: `golcli run -l "##" -i 200 ./buildFile.json`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		buildFilePath = args[0]

		frameInterval, err = cmd.Flags().GetInt("interval")
		livingCellChar, err = cmd.Flags().GetString("live-char")
		deadCellChar, err = cmd.Flags().GetString("dead-char")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		rc := u.RunConfig{
			BuildFilePath: buildFilePath,
			Interval:      time.Duration(frameInterval) * time.Millisecond,
			LiveCellChar:  livingCellChar,
			DeadCellChar:  deadCellChar,
		}

		err = rc.RunBuildFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&frameInterval, "interval", "i", 50, "The number of milliseconds between frames")
	runCmd.Flags().StringVarP(&livingCellChar, "live-char", "l", "  ", "The character(s) that represent a live cell")
	runCmd.Flags().StringVarP(&deadCellChar, "dead-char", "d", "██", "The character(s) that represent a dead cell")
}
