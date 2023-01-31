/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/crodriguezde/rtdashs/pkg/events"
	"github.com/crodriguezde/rtdashs/pkg/generator"
	"github.com/crodriguezde/rtdashs/pkg/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rtdashs",
	Short: "real time dashboards demo",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, err := cmd.Flags().GetString("addr")
		if err != nil {
			return err
		}
		log.Printf("listening on http://%s", addr)

		// Initialize the notify channel
		send := make(chan *events.Cpu)
		ctx := context.Background()

		handler := server.NewHandler(ctx, send)

		// Start generator
		generator.RunCpuGenerator(ctx, send)

		err = http.ListenAndServe(addr, handler)
		if err != nil {
			return err
		}

		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("addr", "a", "localhost:3000", "bind address")
}
