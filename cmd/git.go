package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coreos/go-semver/semver"
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

var gitChangelogCmd = &cobra.Command{
	Use:     string("changelog"),
	Aliases: []string{"c"},
	Short:   string("Generate changelog based on git commits"),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, _ := cmd.Flags().GetString("repo")
		file, _ := cmd.Flags().GetString("file")

		// read file and generate changelog
		content, err := ioutil.ReadFile(filepath.Join(repo, file))
		if err != nil {
			return err
		}
		contentStr := string(content)

		// find the index of the first occurrence of "##"
		index := strings.Index(contentStr, "##")
		eolIndex := strings.Index(contentStr[index:], "\n")
		// match the last semver version
		re := regexp.MustCompile(`\d+\.\d+\.\d+`)
		a := re.FindString(contentStr[index : index+eolIndex])
		fmt.Println(a)

		ver := semver.New(a)
		ver.BumpMinor()
		fmt.Println(ver.String())

		// get the line at index until the end of the line
		line := contentStr[index : index+eolIndex]
		fmt.Println(line)

		// multiline string with new line
		textToInsert := `## [2.303.0] - 2021-03-31

### Added

- Added new feature

### Changed

- Changed something

`

		// insert the new text right before the "##"
		newContentStr := contentStr[:index] + textToInsert + contentStr[index:]

		// get the line at index until the end of the line
		// line := contentStr[:strings.Index(contentStr[index:], "")]
		// fmt.Println(line)

		fmt.Println(newContentStr)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
	gitCmd.AddCommand(gitIgnoreCmd)
	gitCmd.AddCommand(gitChangelogCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	gitChangelogCmd.PersistentFlags().String("repo", ".", "Repository to generate changelog for. Defaults to current directory.")
	gitChangelogCmd.PersistentFlags().String("file", "CHANGELOG.md", "File to write changelog to. Defaults to CHANGELOG.md")

}
