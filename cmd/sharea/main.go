package main

import (
	_ "embed"
	"fmt"
	"github.com/kulikvl/sharea/internal/server"
	"github.com/spf13/cobra"
	"os"
)

func RootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "sharea",
		Short: "Sharea is a fast and easy application for sharing files over LAN",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			folder, _ := cmd.Flags().GetString("folder")
			capacity, _ := cmd.Flags().GetInt64("capacity")

			if port <= 0 || port > 65535 {
				return fmt.Errorf("invalid port number: %d", port)
			}

			if folder != "" {
				if _, err := os.Stat(folder); os.IsNotExist(err) {
					return fmt.Errorf("folder does not exist: %s", folder)
				}
			}

			if capacity <= 0 || capacity > 50*1024*1024*1024 {
				return fmt.Errorf("invalid capacity: %d", capacity)
			}

			serv, err := server.New(port, folder, capacity)
			if err != nil {
				return fmt.Errorf("failed to create a server: %w", err)
			}

			serv.Run()

			return nil
		},
	}

	rootCmd.PersistentFlags().IntP("port", "p", 3000, "Port to run the server on")
	rootCmd.PersistentFlags().StringP("folder", "f", "./storage", "Folder to share")
	rootCmd.PersistentFlags().Int64P("capacity", "c", 1*1024*1024*1024, "Storage capacity in bytes")

	return rootCmd
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
