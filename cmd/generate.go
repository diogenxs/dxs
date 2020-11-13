package cmd

import (
	"errors"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

// Readme projects config
type Readme struct {
	Name string
	// Author string
	// Email  string
	// URL    string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate simple files",
}

// generateReadme generate a README.md file
var generateReadme = &cobra.Command{
	Use:     "readme",
	Short:   "Generate a README.md file",
	Aliases: []string{"r"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name of the project")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpl, err := template.New("readme").Parse(`# {{ .Name }}

## Using

`)
		if err != nil {
			return err
		}

		data := Readme{args[0]}

		outputToFile, _ := cmd.Flags().GetBool("output")
		if !outputToFile {
			err = tmpl.Execute(os.Stdout, data)
			if err != nil {
				return err
			}
			return nil
		}

		output, err := os.Create("README.md")
		if err != nil {
			return err
		}
		defer output.Close()
		err = tmpl.Execute(output, data)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.AddCommand(generateReadme)
	generateReadme.Flags().BoolP("output", "o", false, "Output to a README.md file in current dir")
}
