/*
Copyright © 2025 Rick Ryan
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	addr      string
	serName   string
	env       string
	dbConnStr string
	level     string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().
		StringVarP(&addr, "address", "a", ":8080", "Address of server (host:port)")
	serverCmd.PersistentFlags().StringVarP(&serName, "name", "n", "server", "Server name")
	serverCmd.PersistentFlags().StringVarP(&env, "env", "e", "dev", "Server environment")
	serverCmd.PersistentFlags().
		StringVarP(&dbConnStr, "db-connstr", "c", "postgres://username:password@localhost:5432/database_name", "Connection string for database")
	serverCmd.PersistentFlags().StringVarP(&level, "level", "l", "INFO", "log level of the server")

	err := viper.BindPFlag("server.address", serverCmd.Flags().Lookup("address"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("server.name", serverCmd.Flags().Lookup("name"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("server.env", serverCmd.Flags().Lookup("env"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("server.db-connstr", serverCmd.Flags().Lookup("db-connstr"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("server.level", serverCmd.Flags().Lookup("level"))
	if err != nil {
		panic(err)
	}
}
