package cmd

import (
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// myipCmd represents the myip command
var myipCmd = &cobra.Command{
	Use:   "myip",
	Short: "Show public IP",
	Long:  `Show public IP by consulting https://ifconfig.me`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get("https://ifconfig.me")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		io.Copy(os.Stdout, resp.Body)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(myipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// myipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// myipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
