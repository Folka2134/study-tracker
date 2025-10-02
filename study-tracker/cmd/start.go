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

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [duration] [task]",
	Short: "Starts a new pomodoro timer with a given duration and task.",
	Long:  `Starts a new pomodoro timer with a given duration and task. Once the timer is complete, a desktop notification will be sent to the user.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		duration, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: duration must be an integer.")
			os.Exit(1)
		}
		task := args[1]
		notification.Send("Timer Started!", fmt.Sprintf("Starting a %d minute timer for %s.", duration, task))
		tui.Start(time.Duration(duration)*time.Minute, task, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
