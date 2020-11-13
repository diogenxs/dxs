package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:     "git",
	Short:   "Git helpers",
	Aliases: []string{"g"},
}

// gitCmd represents the git command
var gitIgnoreCmd = &cobra.Command{
	Use:     "ignore",
	Aliases: []string{"i"},
	Short:   "Git helpers",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "https://www.toptal.com/developers/gitignore/api/" + strings.TrimSpace(args[0])
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			fmt.Println("Error in request, response status code: ", resp.StatusCode)
			return nil
		}
		defer resp.Body.Close()

		output, err := os.Create(".gitignore")
		if err != nil {
			return err
		}
		defer output.Close()

		io.Copy(output, resp.Body)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
	gitCmd.AddCommand(gitIgnoreCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
