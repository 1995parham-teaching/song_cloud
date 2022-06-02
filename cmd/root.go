package cmd

import (
	"os"

	"github.com/cng-by-example/song_cloud/cmd/migrate"
	"github.com/cng-by-example/song_cloud/cmd/server"
	"github.com/cng-by-example/song_cloud/config"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.Read()

	// nolint: exhaustruct
	root := &cobra.Command{
		Use:   "song cloud",
		Short: "song cloud",
	}

	server.Register(root, cfg)
	migrate.Register(root, cfg)

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
