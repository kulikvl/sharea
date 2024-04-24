package main

import (
	_ "embed"
	"fmt"
	"github.com/kulikvl/sharea/internal/server"
	"github.com/spf13/cobra"
	"os"
)

func RootCmd() *cobra.Command {
	var port int
	var path string

	var rootCmd = &cobra.Command{
		Use:   "sharea",
		Short: "Sharea is a quick and easy command-line LAN file share utility",
		RunE: func(cmd *cobra.Command, args []string) error {
			if port <= 0 || port > 65535 {
				return fmt.Errorf("invalid port number: %d", port)
			}

			if path != "" {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					return fmt.Errorf("folder does not exist: %s", path)
				}
			}

			serv, err := server.New(port, path, 1000)
			if err != nil {
				return fmt.Errorf("failed to create a server: %w", err)
			}

			serv.Run()

			return nil
		},
	}

	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 3000, "Port to run the server on")
	rootCmd.PersistentFlags().StringVarP(&path, "folder", "f", "./storage", "Folder to share")

	return rootCmd
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
