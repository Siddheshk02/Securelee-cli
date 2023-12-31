/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of Securelee Vault.",
	Long:  `Logout of Securelee Vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := lib.Logout()
		if err != nil {
			fmt.Print("\033[31m", "\n > Error, Please try Again\n", "\033[0m")
			os.Exit(0)
		}
		fmt.Println("\033[36m", "\n > Successfully Logged out!!", "\033[0m")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
