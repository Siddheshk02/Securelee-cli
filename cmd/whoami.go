/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "The Current User of Securelee Vault.",
	Long:  `The Current User of Securelee Vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("whoami called")
		res := lib.Check()
		if !res {
			fmt.Println("\n > No User logged in, You must Login to use Securelee Vault Services.")
			fmt.Println("")
			os.Exit(0)
		}

		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		filePath := filepath.Join(user.HomeDir, "/securelee/token.json")

		jsonData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		var userData struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		err = json.Unmarshal(jsonData, &userData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\n > Logged In as : ", userData.Name, " (", userData.Email, ")")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoamiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoamiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
