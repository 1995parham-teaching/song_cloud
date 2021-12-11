package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"song_cloud/cmd/migrate"
	"song_cloud/cmd/server"
	"song_cloud/config"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.Read()

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "song_cloud",
		Short: "song_cloud",
	}

	server.Register(root, cfg)
	migrate.Register(root, cfg)

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
