package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	QueueName string
)

func init() {
	workerCmd.Flags().StringVar(&QueueName, "queue", "Main", "Queue")
}

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run worker",
	Run: func(cmd *cobra.Command, args []string) {

		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading .env file")
			os.Exit(1)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	// workerCmd.Use = appName

	// Silence Cobra's internal handling of command usage help text.
	// Note, the help text is still displayed if any command arg or
	// flag validation fails.
	workerCmd.SilenceUsage = true

	// Silence Cobra's internal handling of error messaging
	// since we have a custom error handler in main.go
	workerCmd.SilenceErrors = true

	err := workerCmd.Execute()
	return err
}
