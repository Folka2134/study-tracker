package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/folka2134/study-tracker/study-tracker/internal/notification"
	"github.com/folka2134/study-tracker/study-tracker/internal/tui"
	"github.com/spf13/cobra"
)

// breakCmd represents the break command
var breakCmd = &cobra.Command{
	Use:   "break [duration]",
	Short: "Starts a new break timer.",
	Long:  `Starts a new break timer. Break sessions are not saved.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		duration, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: duration must be an integer.")
			os.Exit(1)
		}
		notification.Send("Break Started!", fmt.Sprintf("Starting a %d minute break.", duration))
		tui.Start(time.Duration(duration)*time.Minute, "break", os.Stdout, "#074BF5", "#07F5B9", true)
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)
	// breakCmd.Flags().StringVar(&breakColor, "color", "#00FF00", "Sets the color of the progress bar")
}

