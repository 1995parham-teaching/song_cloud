package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/cng-by-example/song_cloud/cmd/migrate"
	"github.com/cng-by-example/song_cloud/cmd/server"
	"github.com/cng-by-example/song_cloud/config"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.Read()

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "github.com/cng-by-example/song_cloud",
		Short: "github.com/cng-by-example/song_cloud",
	}

	server.Register(root, cfg)
	migrate.Register(root, cfg)

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
