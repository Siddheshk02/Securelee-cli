/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of Securelee Vault.",
	Long:  `Logout of Securelee Vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("logout called")
		err := lib.Logout()
		if err != nil {
			// log.Fatal(err, "Error, Please try again.")
			color.Red("\033[31m", "\n > Error, Please try Again\n", "\033[0m")
			os.Exit(0)
		}
		color.Cyan("\033[36m", "\n > Successfully Logged out!!\n", "\033[0m")
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
