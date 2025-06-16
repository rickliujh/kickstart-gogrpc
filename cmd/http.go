/*
Copyright © 2025 Rick Ryan
*/
package cmd

import (
	server "github.com/rickliujh/kickstart-gogrpc/pkg/server"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "serve using http server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.HTTPServer(addr, serName, env, dbConnStr)
	},
}

func init() {
	serverCmd.AddCommand(httpCmd)
}
