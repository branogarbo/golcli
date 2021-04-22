package cmd

import (
	"fmt"
	"os"

	u "github.com/branogarbo/golcli/util"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Creates a build file from a pattern file.",
	Args:    cobra.ExactArgs(2),
	Example: `golcli build ./pattern.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		patternPath = args[0]
		buildFilePath = args[1]

		frameWidth, err = cmd.Flags().GetInt("width")
		frameHeight, err = cmd.Flags().GetInt("height")
		frameCount, err = cmd.Flags().GetInt("count")
		patternX, err = cmd.Flags().GetInt("pattern-x")
		patternY, err = cmd.Flags().GetInt("pattern-y")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gameConfig := u.GameConfig{
			Width:      frameWidth,
			Height:     frameHeight,
			FrameCount: frameCount,
		}

		initPattern := u.Pattern{
			Path: patternPath,
			X:    patternX,
			Y:    patternY,
		}

		err = u.GenBuildFile(gameConfig, initPattern, buildFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().IntVarP(&frameWidth, "width", "W", 40, "The width of the frames")
	buildCmd.Flags().IntVarP(&frameHeight, "height", "H", 30, "The height of the frames")
	buildCmd.Flags().IntVarP(&frameCount, "count", "c", 10000, "The number of frames displayed before exiting")
	buildCmd.Flags().IntVarP(&patternX, "pattern-x", "x", 12, "The x offset of the initial pattern")
	buildCmd.Flags().IntVarP(&patternY, "pattern-y", "y", 8, "The y offset of the initial pattern")
}
