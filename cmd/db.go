/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/diogenxs/dxs/models"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "run commands against local database",
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "run migrations on local database",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("migrating database")
		err := models.MigrateDB()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(migrateCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
