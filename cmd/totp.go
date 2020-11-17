package cmd

import (
	"fmt"

	totp "github.com/hgfischer/go-otp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// totpCmd represents the totp command
var totpCmd = &cobra.Command{
	Use:     "totp",
	Aliases: []string{"t"},
	Short:   "Generate TOTP codes.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m := viper.GetStringMapString("totp")
		secret := m[args[0]]
		if secret == "" {
			fmt.Println("Key not found")
			return
		}

		t := &totp.TOTP{Secret: secret, IsBase32Secret: true}
		fmt.Println(t.Get())
	},
}

func init() {
	rootCmd.AddCommand(totpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// totpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// totpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("viper", rootCmd.PersistentFlags().Lookup("viper"))
}
