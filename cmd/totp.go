package cmd

import (
	"fmt"
	"log"

	"github.com/diogenxs/dxs/utils"
	totp "github.com/hgfischer/go-otp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// totpCmd represents the totp command
var totpCmd = &cobra.Command{
	Use:     "totp",
	Aliases: []string{"t"},
	Short:   "Generate TOTP codes.",
	// Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m := viper.GetStringMapString("totp")
		if len(m) == 0 {
			fmt.Println("No keys added to config file: ", viper.ConfigFileUsed())
		}
		if len(args) < 1 {
			fmt.Println("Possible keys:")
			for k := range m {
				fmt.Println(k)
			}
		} else {
			secret := m[args[0]]
			if viper.GetBool("verbose") {
				fmt.Println("getting key:", args[0])
			}
			if secret == "" {
				log.Fatal("key not found")
			}

			t := &totp.TOTP{Secret: secret, IsBase32Secret: true}
			token := t.Get()

			if v, _ := cmd.Flags().GetBool("copy"); v {
				utils.NewNotification("dxs", fmt.Sprintf("TOTP \"%s\" copied to clipboard", args[0]))
				utils.WriteToClipboard(token)
			} else {
				fmt.Printf("%s", token)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(totpCmd)

	totpCmd.Flags().BoolP("copy", "c", false, "Copy the secret to clipboard")
}
