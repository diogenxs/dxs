/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/diogenxs/dxs/models"
	"github.com/diogenxs/dxs/pkg"
	"github.com/spf13/cobra"
)

// alertsCmd represents the alerts command
var alertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "A brief description of your command",
}
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := models.GetDB()
		if err != nil {
			return err
		}
		defer db.Close()

		alerts, err := models.ListAlerts(db)
		if err != nil {
			return err
		}
		for _, alert := range alerts {
			fmt.Println(alert.Fingerprint + ": " + alert.Labels["alertname"])
		}
		return nil
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync new alerts from third-party apps to the local database",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := pkg.SyncAlerts()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(alertsCmd)

	alertsCmd.AddCommand(listCmd)
	alertsCmd.AddCommand(syncCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
