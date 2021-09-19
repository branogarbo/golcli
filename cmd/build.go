package cmd

import (
	u "github.com/branogarbo/golcli/util"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Creates a build file from a pattern file.",
	Example: `golcli build -W 70 -c 400 pattern.txt builds/pattern`,
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var buildFilePath string
		patternFilePath := args[0]

		if len(args) == 2 {
			buildFilePath = args[1]
		}

		frameWidth, _ := cmd.Flags().GetInt("width")
		frameHeight, _ := cmd.Flags().GetInt("height")
		frameCount, _ := cmd.Flags().GetInt("count")
		patternX, _ := cmd.Flags().GetInt("pattern-x")
		patternY, _ := cmd.Flags().GetInt("pattern-y")

		bc := u.BuildConfig{
			BuildFilePath: buildFilePath,
			Width:         frameWidth,
			Height:        frameHeight,
			FrameCount:    frameCount,
			InitPattern: u.Pattern{
				FilePath: patternFilePath,
				X:        patternX,
				Y:        patternY,
			},
		}

		return u.GenBuildFile(bc)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().IntP("width", "W", 40, "The width of the frames")
	buildCmd.Flags().IntP("height", "H", 30, "The height of the frames")
	buildCmd.Flags().IntP("count", "c", 500, "The number of frames displayed before exiting")
	buildCmd.Flags().IntP("pattern-x", "x", 12, "The x offset of the initial pattern")
	buildCmd.Flags().IntP("pattern-y", "y", 8, "The y offset of the initial pattern")
}
