package cmd

import (
	"time"

	u "github.com/branogarbo/golcli/util"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Runs a build file.",
	Example: `golcli run -l "##" -d ".." -i 200 buildFiles/pattern`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		buildFilePath := args[0]

		frameInterval, _ := cmd.Flags().GetInt("interval")
		livingCellChar, _ := cmd.Flags().GetString("live-char")
		deadCellChar, _ := cmd.Flags().GetString("dead-char")

		rc := u.RunConfig{
			BuildFilePath: buildFilePath,
			Interval:      time.Duration(frameInterval) * time.Millisecond,
			LiveCellChar:  livingCellChar,
			DeadCellChar:  deadCellChar,
		}

		return u.RunBuildFile(rc)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntP("interval", "i", 50, "The number of milliseconds between frames")
	runCmd.Flags().StringP("live-char", "l", "  ", "The character(s) that represent a live cell")
	runCmd.Flags().StringP("dead-char", "d", "██", "The character(s) that represent a dead cell")
}
